package build

import (
	"context"
	"net"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/pkg/middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func (b *Builder) gRPCServer(ctx context.Context) *grpc.Server {
	logger := zerolog.Ctx(ctx)
	//unaryOpts := []grpc.UnaryServerInterceptor{
	//	middleware.GRPCRecoverUnaryServerInterceptor(logger, true),
	//	middleware.GRPCAccessLogUnaryServerInterceptor(logger),
	//}

	streamOpts := []grpc.StreamServerInterceptor{
		middleware.GRPCRecoverStreamServerInterceptor(logger, true),
		middleware.GRPCAccessLogStreamServerInterceptor(logger),
	}

	serverOpts := make([]grpc.ServerOption, 0)
	serverOpts = append(
		serverOpts,
		//grpc.ChainUnaryInterceptor(unaryOpts...),
		grpc.ChainStreamInterceptor(streamOpts...),
	)

	grpcServer := grpc.NewServer(serverOpts...)

	b.shutdown.add(func(ctx context.Context) error {
		grpcServer.Stop()

		return nil
	})

	return grpcServer
}

func (b *Builder) Listener() (net.Listener, error) {
	listener, err := net.Listen(b.config.GRPC.NetworkType, b.config.GRPCAddr())
	if err != nil {
		return nil, errors.Wrap(err, "start network listener")
	}

	b.logger.Info().Str("address", b.config.GRPCAddr()).Msg("Listen gRPC")

	return listener, nil
}
