package pdtesting

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/fetcher"
	"github.com/weaveworks/weave-gitops/core/nsaccess/nsaccessfakes"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
	v1 "k8s.io/api/core/v1"
	v1a "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

func MakeHTTPServer(
	t *testing.T,
	k8sEnv *testutils.K8sTestEnv,
) *httptest.Server {
	log := logr.Discard()
	ctx := context.Background()

	cl, err := RestConfigToCluster(k8sEnv.Rest)
	if err != nil {
		t.Fatal(err)
	}

	fetcher := fetcher.NewSingleClusterFetcher(cl)

	nsChecker := nsaccessfakes.FakeChecker{}
	nsChecker.FilterAccessibleNamespacesStub = func(ctx context.Context, _ v1a.AuthorizationV1Interface, n []v1.Namespace) ([]v1.Namespace, error) {
		// Pretend the user has access to everything
		return n, nil
	}

	clustersManager := clustersmngr.NewClustersManager([]clustersmngr.ClusterFetcher{fetcher}, &nsChecker, log)

	_ = clustersManager.UpdateClusters(ctx)
	_ = clustersManager.UpdateNamespaces(ctx)

	opts := server.ServerOpts{
		ClustersManager: clustersManager,
		CRDService:      crd.NewNoCacheFetcher(clustersManager),
	}

	mux := runtime.NewServeMux()

	if err := server.Hydrate(ctx, mux, opts); err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(mux)

	t.Cleanup(func() {
		ts.Close()
	})

	return ts
}
