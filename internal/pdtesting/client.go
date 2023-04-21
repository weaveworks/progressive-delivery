package pdtesting

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/fetcher"
	"github.com/weaveworks/weave-gitops/core/nsaccess"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
)

func CreateClient(k8sEnv *testutils.K8sTestEnv) (clustersmngr.Client, clustersmngr.ClustersManager, error) {
	ctx := context.Background()
	log := logr.Discard()

	cl, err := RestConfigToCluster(k8sEnv.Rest)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create cluster from config: %w", err)
	}

	fetcher := fetcher.NewSingleClusterFetcher(cl)

	clustersManager := clustersmngr.NewClustersManager(
		[]clustersmngr.ClusterFetcher{fetcher},
		nsaccess.NewChecker(nsaccess.DefautltWegoAppRules),
		log,
	)

	if err := clustersManager.UpdateClusters(ctx); err != nil {
		return nil, nil, err
	}
	if err := clustersManager.UpdateNamespaces(ctx); err != nil {
		return nil, nil, err
	}

	client, err := clustersManager.GetServerClient(ctx)

	return client, clustersManager, err
}
