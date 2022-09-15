package server

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	"sigs.k8s.io/kustomize/kstatus/status"
)

func (pd *pdServer) ListCanaryObjects(ctx context.Context, msg *pb.ListCanaryObjectsRequest) (*pb.ListCanaryObjectsResponse, error) {
	clusterClient, err := pd.clustersManager.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonating client: %w", err)
	}

	respErrors := []*pb.ListError{}
	result, err := pd.flagger.ListCanaryObjects(ctx, clusterClient, flagger.ListCanaryObjectsOptions{
		Name:        msg.Name,
		Namespace:   msg.Namespace,
		ClusterName: msg.ClusterName,
	})
	if err != nil {
		if merr, ok := err.(*multierror.Error); ok {
			for _, err := range merr.Errors {
				respErrors = append(respErrors, &pb.ListError{ClusterName: msg.ClusterName, Namespace: msg.Namespace, Message: err.Error()})
			}
		} else {
			return nil, fmt.Errorf("failed listing canary objects: %w", err)
		}
	}

	objects := []*pb.UnstructuredObject{}

	for _, obj := range result {
		res, err := status.Compute(&obj)
		if err != nil {
			return nil, fmt.Errorf("could not get status for %s: %w", obj.GetName(), err)
		}

		var images []string

		switch obj.GetKind() {
		case "Deployment":
			images = getDeploymentPodContainerImages(obj.Object)
		}

		objects = append(objects, &pb.UnstructuredObject{
			GroupVersionKind: &pb.GroupVersionKind{
				Group:   obj.GetObjectKind().GroupVersionKind().Group,
				Version: obj.GetObjectKind().GroupVersionKind().GroupVersion().Version,
				Kind:    obj.GetKind(),
			},
			Name:        obj.GetName(),
			Namespace:   obj.GetNamespace(),
			Images:      images,
			Status:      res.Status.String(),
			Uid:         string(obj.GetUID()),
			Conditions:  mapUnstructuredConditions(res),
			ClusterName: msg.GetClusterName(),
		})
	}

	return &pb.ListCanaryObjectsResponse{Objects: objects, Errors: respErrors}, nil
}
