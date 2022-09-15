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
	pds, err := NewProgressiveDeliveryServer(opts)
	if err != nil {
		return err
	}

	return pb.RegisterProgressiveDeliveryServiceHandlerServer(ctx, mux, pds)
}

type pdServer struct {
	pb.UnimplementedProgressiveDeliveryServiceServer

	clustersManager clustersmngr.ClustersManager
	version         version.Fetcher
	crd             crd.Fetcher
	flagger         flagger.Fetcher
	logger          logr.Logger
}

type ServerOpts struct {
	ClustersManager clustersmngr.ClustersManager
	CRDService      crd.Fetcher
	Logger          logr.Logger
}

func NewProgressiveDeliveryServer(opts ServerOpts) (pb.ProgressiveDeliveryServiceServer, error) {
	ctx := context.Background()

	versionService := version.NewFetcher()

	if opts.CRDService == nil {
		opts.CRDService = crd.NewFetcher(ctx, opts.Logger, opts.ClustersManager)
	}

	flaggerService := flagger.NewFetcher(opts.CRDService, opts.Logger)

	return &pdServer{
		clustersManager: opts.ClustersManager,
		version:         versionService,
		crd:             opts.CRDService,
		flagger:         flaggerService,
		logger:          opts.Logger,
	}, nil
}

func (pd *pdServer) GetVersion(ctx context.Context, msg *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	v := pd.version.Get()

	return &pb.GetVersionResponse{Version: v.String()}, nil
}
