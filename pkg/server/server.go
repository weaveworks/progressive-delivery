package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/services/version"
)

func Hydrate(ctx context.Context, mux *runtime.ServeMux, opts ServerOpts) error {
	pds := NewProgressiveDeliveryServer(opts)
	return pb.RegisterProgressiveDeliveryServiceHandlerServer(ctx, mux, pds)
}

type pdServer struct {
	pb.UnimplementedProgressiveDeliveryServiceServer

	version version.Fetcher
}

type ServerOpts struct{}

func NewProgressiveDeliveryServer(opts ServerOpts) pb.ProgressiveDeliveryServiceServer {
	return &pdServer{
		version: version.NewFetcher(),
	}
}

func (pd *pdServer) GetVersion(ctx context.Context, msg *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	v := pd.version.Get()

	return &pb.GetVersionResponse{Version: v.String()}, nil
}
