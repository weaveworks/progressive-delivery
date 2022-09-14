package pdtesting

import (
	"context"
	"net"
	"testing"

	"github.com/go-logr/logr"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/clustersmngrfakes"
	"github.com/weaveworks/weave-gitops/core/nsaccess/nsaccessfakes"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

func MakeGRPCServer(
	t *testing.T,
	cfg *rest.Config,
	k8sEnv *testutils.K8sTestEnv,
) pb.ProgressiveDeliveryServiceClient {
	log := logr.Discard()
	ctx := context.Background()

	fetcher := &clustersmngrfakes.FakeClusterFetcher{}
	fetcher.FetchReturns([]clustersmngr.Cluster{RestConfigToCluster(k8sEnv.Rest)}, nil)

	nsChecker := nsaccessfakes.FakeChecker{}
	nsChecker.FilterAccessibleNamespacesStub = func(ctx context.Context, c *rest.Config, n []v1.Namespace) ([]v1.Namespace, error) {
		// Pretend the user has access to everything
		return n, nil
	}

	clientsFactory := clustersmngr.NewClustersManager(
		fetcher,
		&nsChecker,
		log,
		kube.CreateScheme(),
		clustersmngr.NewClustersClientsPool,
		clustersmngr.DefaultKubeConfigOptions,
	)

	_ = clientsFactory.UpdateClusters(ctx)
	_ = clientsFactory.UpdateNamespaces(ctx)

	opts := server.ServerOpts{
		ClientFactory: clientsFactory,
		CRDService:    crd.NewNoCacheFetcher(clientsFactory),
		Logger:        logr.Discard(),
	}

	pdServer, _ := server.NewProgressiveDeliveryServer(opts)
	lis := bufconn.Listen(1024 * 1024)
	principal := auth.NewUserPrincipal(auth.Token("1234"))
	s := grpc.NewServer(
		withClientsPoolInterceptor(clientsFactory, cfg, principal),
	)

	pb.RegisterProgressiveDeliveryServiceServer(s, pdServer)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	go func(tt *testing.T) {
		if err := s.Serve(lis); err != nil {
			tt.Error(err)
		}
	}(t)

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		s.GracefulStop()
		conn.Close()
	})

	return pb.NewProgressiveDeliveryServiceClient(conn)
}

func RestConfigToCluster(cfg *rest.Config) clustersmngr.Cluster {
	return clustersmngr.Cluster{
		Name:        "Default",
		Server:      cfg.Host,
		BearerToken: cfg.BearerToken,
		TLSConfig:   cfg.TLSClientConfig,
	}
}

func withClientsPoolInterceptor(clientsFactory clustersmngr.ClustersManager, config *rest.Config, user *auth.UserPrincipal) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := clientsFactory.UpdateClusters(ctx); err != nil {
			return nil, err
		}
		if err := clientsFactory.UpdateNamespaces(ctx); err != nil {
			return nil, err
		}

		clientsFactory.UpdateUserNamespaces(ctx, user)

		ctx = auth.WithPrincipal(ctx, user)

		clusterClient, err := clientsFactory.GetImpersonatedClient(ctx, user)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, clustersmngr.ClustersClientCtxKey, clusterClient)

		return handler(ctx, req)
	})
}
