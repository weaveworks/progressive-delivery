package server

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/progressive-delivery/pkg/services/version"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
)

func Hydrate(ctx context.Context, mux *runtime.ServeMux, opts ServerOpts) error {
	pds, err := NewProgressiveDeliveryServer(ctx, opts)
	if err != nil {
		return err
	}

	return pb.RegisterProgressiveDeliveryServiceHandlerServer(ctx, mux, pds)
}

type pdServer struct {
	pb.UnimplementedProgressiveDeliveryServiceServer

	logger         logr.Logger
	clientsFactory clustersmngr.ClientsFactory
	version        version.Fetcher
	crd            crd.Fetcher
	flagger        flagger.Fetcher
}

type ServerOpts struct {
	Logger        logr.Logger
	ClientFactory clustersmngr.ClientsFactory
	CRDService    crd.Fetcher
}

func NewProgressiveDeliveryServer(ctx context.Context, opts ServerOpts) (pb.ProgressiveDeliveryServiceServer, error) {
	logger := opts.Logger.WithName("progressive-delivery-server")
	clusterClient, err := opts.ClientFactory.GetServerClient(ctx)
	if err != nil {
		logger.Error(err, "failed to get server client")
		return nil, err
	}
	logger.Info("listing clients", "clients", clusterClient.ClientsPool().Clients())

	versionService := version.NewFetcher()

	if opts.CRDService == nil {
		opts.CRDService = crd.NewFetcher(ctx, opts.Logger, clusterClient)
	}

	flaggerService := flagger.NewFetcher(opts.CRDService)

	return &pdServer{
		logger:         logger,
		clientsFactory: opts.ClientFactory,
		version:        versionService,
		crd:            opts.CRDService,
		flagger:        flaggerService,
	}, nil
}

func (pd *pdServer) GetVersion(ctx context.Context, msg *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	v := pd.version.Get()

	return &pb.GetVersionResponse{Version: v.String()}, nil
}
