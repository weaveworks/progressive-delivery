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

	clusterClient, err := pd.clientsFactory.GetServerClient(ctx)
	if err != nil {
		pd.logger.Error(err, "failed to get server client")
		return nil, err
	}
	pd.logger.Info("listing clients", "clients", clusterClient.ClientsPool().Clients())

	return &pb.IsFlaggerAvailableResponse{
		Clusters: pd.crd.IsAvailableOnClusters(crd.FlaggerCRDName),
	}, nil
}

func (pd *pdServer) ListCanaries(ctx context.Context, msg *pb.ListCanariesRequest) (*pb.ListCanariesResponse, error) {
	clusterClient, err := pd.clientsFactory.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonated client: %w", err)
	}

	results, nextPageToken, listErr, err := pd.flagger.ListCanaryDeployments(
		ctx,
		clusterClient,
		flagger.ListCanaryDeploymentsOptions{},
	)
	if err != nil {
		pd.logger.Error(err, "unable to list canaries", "clients", clusterClient.ClientsPool().Clients())
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

			pbObject := convert.FlaggerCanaryToProto(item, clusterName, deployment)

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
		return nil, err
	}

	deployment, err := pd.flagger.FetchTargetRef(ctx, msg.ClusterName, clusterClient, canary)
	if err != nil {
		return nil, err
	}

	pbObject := convert.FlaggerCanaryToProto(*canary, msg.ClusterName, deployment)

	pbObject.DeploymentStrategy = string(pd.flagger.DeploymentStrategyFor(*canary))

	response := &pb.GetCanaryResponse{
		Canary:     pbObject,
		Automation: getAutomation(deployment),
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
