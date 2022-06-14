package pdtesting

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/clustersmngrfakes"
	"github.com/weaveworks/weave-gitops/core/nsaccess"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
)

func CreateClient(k8sEnv *testutils.K8sTestEnv) (clustersmngr.Client, clustersmngr.ClientsFactory, error) {
	ctx := context.Background()
	log := logr.Discard()
	fetcher := &clustersmngrfakes.FakeClusterFetcher{}
	fetcher.FetchReturns([]clustersmngr.Cluster{RestConfigToCluster(k8sEnv.Rest)}, nil)

	clientsFactory := clustersmngr.NewClientFactory(
		fetcher,
		nsaccess.NewChecker(nsaccess.DefautltWegoAppRules),
		log,
		server.CreateScheme(),
	)

	if err := clientsFactory.UpdateClusters(ctx); err != nil {
		return nil, nil, err
	}
	if err := clientsFactory.UpdateNamespaces(ctx); err != nil {
		return nil, nil, err
	}

	client, err := clientsFactory.GetServerClient(ctx)

	return client, clientsFactory, err
}
