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

	"github.com/urfave/cli/v2"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/core/clustersmngr/fetcher"
	"github.com/weaveworks/weave-gitops/core/logger"
	"github.com/weaveworks/weave-gitops/core/nsaccess/nsaccessfakes"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type appConfig struct {
	Host string
	Port string
}

func NewApp(out io.Writer) *cli.App {
	cfg := &appConfig{}

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
	log, err := logger.New("debug", true)
	if err != nil {
		return err
	}

	ctx := context.Background()

	restCfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not create client config: %w", err)
	}

	fetcher := fetcher.NewSingleClusterFetcher(restCfg)

	nsChecker := nsaccessfakes.FakeChecker{}
	nsChecker.FilterAccessibleNamespacesStub = func(ctx context.Context, c *rest.Config, n []v1.Namespace) ([]v1.Namespace, error) {
		// Pretend the user has access to everything
		return n, nil
	}

	clientsFactory := clustersmngr.NewClientFactory(
		fetcher,
		&nsChecker,
		log,
		server.CreateScheme(),
	)
	clientsFactory.Start(ctx)

	_ = clientsFactory.UpdateClusters(ctx)
	_ = clientsFactory.UpdateNamespaces(ctx)

	clusterClient, err := clientsFactory.GetServerClient(ctx)
	if err != nil {
		return err
	}

	opts := server.ServerOpts{
		ClientFactory: clientsFactory,
		CRDService:    crd.NewNoCacheFetcher(clusterClient),
	}

	pdServer, _ := server.NewProgressiveDeliveryServer(opts)
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	principal := &auth.UserPrincipal{}
	s := grpc.NewServer(
		withClientsPoolInterceptor(clientsFactory, restCfg, principal),
	)

	pb.RegisterProgressiveDeliveryServiceServer(s, pdServer)

	reflection.Register(s)

	go func() {
		log.Info("Starting server", "address", address)

		if err := s.Serve(lis); err != nil {
			log.Error(err, "server exited")
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

func withClientsPoolInterceptor(clientsFactory clustersmngr.ClientsFactory, config *rest.Config, user *auth.UserPrincipal) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := clientsFactory.UpdateClusters(ctx); err != nil {
			return nil, err
		}
		if err := clientsFactory.UpdateNamespaces(ctx); err != nil {
			return nil, err
		}

		clientsFactory.UpdateUserNamespaces(ctx, user)

		clusterClient, err := clientsFactory.GetImpersonatedClient(ctx, user)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, clustersmngr.ClustersClientCtxKey, clusterClient)

		return handler(ctx, req)
	})
}
