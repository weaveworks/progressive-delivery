package server

import (
	"context"

	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/convert"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
)

func (pd *pdServer) ListCanaries(ctx context.Context, msg *pb.ListCanariesRequest) (*pb.ListCanariesResponse, error) {
	clusterClient := clustersmngr.ClientFromCtx(ctx)

	results, nextPageToken, listErr, err := pd.flagger.ListCanaryDeployments(
		ctx,
		clusterClient,
		flagger.ListCanaryDeploymentsOptions{},
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

			pbObject := convert.FlaggerCanaryToProto(item, clusterName, deployment)
			response.Canaries = append(response.Canaries, pbObject)
		}
	}

	return response, nil
}
