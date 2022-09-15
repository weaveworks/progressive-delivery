package pdtesting

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/clustersmngrfakes"
	"github.com/weaveworks/weave-gitops/core/nsaccess"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
)

func CreateClient(k8sEnv *testutils.K8sTestEnv) (clustersmngr.Client, clustersmngr.ClustersManager, error) {
	ctx := context.Background()
	log := logr.Discard()
	fetcher := &clustersmngrfakes.FakeClusterFetcher{}
	fetcher.FetchReturns([]clustersmngr.Cluster{RestConfigToCluster(k8sEnv.Rest)}, nil)

	clustersManager := clustersmngr.NewClustersManager(
		fetcher,
		nsaccess.NewChecker(nsaccess.DefautltWegoAppRules),
		log,
		kube.CreateScheme(),
		clustersmngr.NewClustersClientsPool,
		clustersmngr.DefaultKubeConfigOptions,
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
