package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	stdlog "log"

	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/cluster"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/fetcher"
	"github.com/weaveworks/weave-gitops/core/logger"
	"github.com/weaveworks/weave-gitops/core/nsaccess/nsaccessfakes"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/api/core/v1"
	v1a "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type appConfig struct {
	Host   string
	Port   string
	Logger logr.Logger
}

func NewApp(out io.Writer) *cli.App {
	log, err := logger.New(logger.DefaultLogLevel, os.Getenv("HUMAN_LOGS") != "")
	if err != nil {
		stdlog.Fatalf("Couldn't set up logger: %v", err)
	}

	cfg := &appConfig{
		Logger: log,
	}

	app := &cli.App{
		Name:  "server",
		Usage: "Progressive Delivery Server",
		Flags: CLIFlags(
			WithHTTPServerFlags(),
		),
		Before: parseFlags(cfg),
		Action: func(c *cli.Context) error {
			return serve(cfg)
		},
	}

	if out != nil {
		app.Writer = out
	}

	return app
}

func serve(cfg *appConfig) error {
	ctx := context.Background()

	restCfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not create client config: %w", err)
	}

	scheme := kube.CreateScheme()

	cl, err := cluster.NewSingleCluster(cluster.DefaultCluster, restCfg, scheme, cluster.DefaultKubeConfigOptions...)
	if err != nil {
		return fmt.Errorf("unable to create single cluster: %w", err)
	}

	fetcher := fetcher.NewSingleClusterFetcher(cl)

	nsChecker := nsaccessfakes.FakeChecker{}
	nsChecker.FilterAccessibleNamespacesStub = func(ctx context.Context, _ v1a.AuthorizationV1Interface, n []v1.Namespace) ([]v1.Namespace, error) {
		// Pretend the user has access to everything
		return n, nil
	}

	clustersManager := clustersmngr.NewClustersManager([]clustersmngr.ClusterFetcher{fetcher}, &nsChecker, cfg.Logger)
	clustersManager.Start(ctx)

	_ = clustersManager.UpdateClusters(ctx)
	_ = clustersManager.UpdateNamespaces(ctx)

	opts := server.ServerOpts{
		ClustersManager: clustersManager,
		CRDService:      crd.NewNoCacheFetcher(clustersManager),
		Logger:          cfg.Logger,
	}

	pdServer, _ := server.NewProgressiveDeliveryServer(opts)
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	principal := &auth.UserPrincipal{
		ID:     "pd-admin",
		Groups: []string{"admin"},
	}
	s := grpc.NewServer(
		withClientsPoolInterceptor(clustersManager, restCfg, principal),
	)

	pb.RegisterProgressiveDeliveryServiceServer(s, pdServer)

	reflection.Register(s)

	go func() {
		cfg.Logger.Info("Starting server", "address", address)

		if err := s.Serve(lis); err != nil {
			cfg.Logger.Error(err, "server exited")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer func() {
		cancel()
	}()

	s.GracefulStop()

	return nil
}

func withClientsPoolInterceptor(clustersManager clustersmngr.ClustersManager, config *rest.Config, user *auth.UserPrincipal) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := clustersManager.UpdateClusters(ctx); err != nil {
			return nil, err
		}
		if err := clustersManager.UpdateNamespaces(ctx); err != nil {
			return nil, err
		}

		clustersManager.UpdateUserNamespaces(ctx, user)

		clusterClient, err := clustersManager.GetImpersonatedClient(ctx, user)
		if err != nil {
			return nil, err
		}

		ctx = auth.WithPrincipal(ctx, user)
		ctx = context.WithValue(ctx, clustersmngr.ClustersClientCtxKey, clusterClient)

		return handler(ctx, req)
	})
}
