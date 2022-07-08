package pdtesting

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/clustersmngrfakes"
	"github.com/weaveworks/weave-gitops/core/nsaccess/nsaccessfakes"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

func MakeHTTPServer(
	t *testing.T,
	k8sEnv *testutils.K8sTestEnv,
) *httptest.Server {
	log := logr.Discard()
	ctx := context.Background()

	fetcher := &clustersmngrfakes.FakeClusterFetcher{}
	fetcher.FetchReturns([]clustersmngr.Cluster{RestConfigToCluster(k8sEnv.Rest)}, nil)

	nsChecker := nsaccessfakes.FakeChecker{}
	nsChecker.FilterAccessibleNamespacesStub = func(ctx context.Context, c *rest.Config, n []v1.Namespace) ([]v1.Namespace, error) {
		// Pretend the user has access to everything
		return n, nil
	}

	clientsFactory := clustersmngr.NewClientFactory(
		fetcher,
		&nsChecker,
		log,
		kube.CreateScheme(),
	)

	_ = clientsFactory.UpdateClusters(ctx)
	_ = clientsFactory.UpdateNamespaces(ctx)

	opts := server.ServerOpts{
		ClientFactory: clientsFactory,
		CRDService:    crd.NewNoCacheFetcher(clientsFactory),
	}

	mux := runtime.NewServeMux()

	err := server.Hydrate(ctx, mux, opts)
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(mux)

	t.Cleanup(func() {
		ts.Close()
	})

	return ts
}
