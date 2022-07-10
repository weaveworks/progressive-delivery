package server

import (
	"context"
	"fmt"

	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/convert"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	v1 "k8s.io/api/apps/v1"
)

const (
	LabelKustomizeName        = "kustomize.toolkit.fluxcd.io/name"
	LabelKustomizeNamespace   = "kustomize.toolkit.fluxcd.io/namespace"
	LabelHelmReleaseName      = "helm.toolkit.fluxcd.io/name"
	LabelHelmReleaseNamespace = "helm.toolkit.fluxcd.io/namespace"
)

func (pd *pdServer) IsFlaggerAvailable(ctx context.Context, msg *pb.IsFlaggerAvailableRequest) (*pb.IsFlaggerAvailableResponse, error) {
	return &pb.IsFlaggerAvailableResponse{
		Clusters: pd.crd.IsAvailableOnClusters(crd.FlaggerCRDName),
	}, nil
}

func (pd *pdServer) ListCanaries(ctx context.Context, msg *pb.ListCanariesRequest) (*pb.ListCanariesResponse, error) {
	clusterClient, err := pd.clientsFactory.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonated client: %w", err)
	}

	opts := flagger.ListCanaryDeploymentsOptions{}
	if msg.Pagination != nil {
		opts.PageSize = msg.Pagination.PageSize
		opts.PageToken = msg.Pagination.PageToken
	}

	results, nextPageToken, listErr, err := pd.flagger.ListCanaryDeployments(
		ctx,
		clusterClient,
		opts,
	)
	if err != nil {
		return nil, err
	}

	response := &pb.ListCanariesResponse{
		Canaries:      []*pb.Canary{},
		NextPageToken: nextPageToken,
		Errors:        []*pb.ListError{},
	}

	for _, err := range listErr {
		response.Errors = append(response.Errors, &pb.ListError{
			ClusterName: err.ClusterName,
			Namespace:   "",
			Message:     err.Error(),
		})
	}

	for clusterName, list := range results {
		for _, item := range list {
			// Ignored intentioannly. The function returns with an error, but here we
			// don't care about it, if it's not found, we can return to the client
			// with an empty deployment.
			deployment, _ := pd.flagger.FetchTargetRef(ctx, clusterName, clusterClient, &item)

			promoted, _ := pd.flagger.FetchPromoted(ctx, clusterName, clusterClient, &item)

			containers := promoted.Spec.Template.Spec.Containers
			pbObject := convert.FlaggerCanaryToProto(item, clusterName, deployment, containers)

			pbObject.DeploymentStrategy = string(pd.flagger.DeploymentStrategyFor(item))

			response.Canaries = append(response.Canaries, pbObject)
		}
	}

	return response, nil
}

func (pd *pdServer) GetCanary(ctx context.Context, msg *pb.GetCanaryRequest) (*pb.GetCanaryResponse, error) {
	clusterClient, err := pd.clientsFactory.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonated client: %w", err)
	}

	canary, err := pd.flagger.GetCanary(ctx, clusterClient, flagger.GetCanaryOptions{
		Name:        msg.Name,
		Namespace:   msg.Namespace,
		ClusterName: msg.ClusterName,
	})
	if err != nil {
		return nil, fmt.Errorf("getting canary: %w", err)
	}

	deployment, err := pd.flagger.FetchTargetRef(ctx, msg.ClusterName, clusterClient, canary)
	if err != nil {
		return nil, fmt.Errorf("fetching target ref: %w", err)
	}

	promoted, err := pd.flagger.FetchPromoted(ctx, msg.ClusterName, clusterClient, canary)
	if err != nil {
		return nil, fmt.Errorf("fetching target ref: %w", err)
	}

	containers := promoted.Spec.Template.Spec.Containers
	pbObject := convert.FlaggerCanaryToProto(*canary, msg.ClusterName, deployment, containers)

	pbObject.DeploymentStrategy = string(pd.flagger.DeploymentStrategyFor(*canary))
	//pbObject.Analysis.Metrics = []*pb.CanaryMetric{}
	//TODO: resolve metric template references before returning metrics
	//for _, item := range canary.GetAnalysis().Metrics {
	//	if item.TemplateRef != nil {
	//		template, err := pd.flagger.GetMetricTemplate(
	//			ctx,
	//			msg.ClusterName,
	//			clusterClient,
	//			item.TemplateRef.Name,
	//			item.TemplateRef.Namespace,
	//		)
	//		if err != nil {
	//			pd.logger.Error(err, "unable to fetch metric template from reference")
	//			continue
	//		}
	//
	//		pbObject.Analysis.Metrics = append(
	//			pbObject.Analysis.Metrics,
	//			convert.FlaggerMetricTemplateToProto(template, msg.ClusterName),
	//		)
	//	}
	//}

	response := &pb.GetCanaryResponse{
		Canary:     pbObject,
		Automation: getAutomation(deployment),
	}

	return response, nil
}

func (pd *pdServer) ListMetricTemplates(ctx context.Context, msg *pb.ListMetricTemplatesRequest) (*pb.ListMetricTemplatesResponse, error) {
	clusterClient, err := pd.clientsFactory.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonated client: %w", err)
	}

	opts := flagger.ListMetricTemplatesOptions{}
	if msg.Pagination != nil {
		opts.PageSize = msg.Pagination.PageSize
		opts.PageToken = msg.Pagination.PageToken
	}

	results, nextPageToken, listErr, err := pd.flagger.ListMetricTemplates(
		ctx,
		clusterClient,
		opts,
	)
	if err != nil {
		return nil, err
	}

	response := &pb.ListMetricTemplatesResponse{
		Templates:     []*pb.CanaryMetricTemplate{},
		NextPageToken: nextPageToken,
		Errors:        []*pb.ListError{},
	}

	for _, err := range listErr {
		response.Errors = append(response.Errors, &pb.ListError{
			ClusterName: err.ClusterName,
			Namespace:   "",
			Message:     err.Error(),
		})
	}

	for clusterName, list := range results {
		for _, item := range list {
			pbObject := convert.FlaggerMetricTemplateToProto(item, clusterName)
			response.Templates = append(response.Templates, pbObject)
		}
	}

	return response, nil
}

func getAutomation(dpl v1.Deployment) *pb.Automation {
	for k, v := range dpl.Labels {
		switch k {
		case LabelKustomizeName:
			return &pb.Automation{
				Kind:      "Kustomization",
				Name:      v,
				Namespace: dpl.Labels[LabelKustomizeNamespace],
			}
		case LabelHelmReleaseName:
			return &pb.Automation{
				Kind:      "HelmRelease",
				Name:      v,
				Namespace: dpl.Labels[LabelHelmReleaseNamespace],
			}
		}
	}

	return nil
}
