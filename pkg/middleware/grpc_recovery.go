package middleware

import (
	"context"

	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCRecoverUnaryServerInterceptor(logger *zerolog.Logger, enableStack bool) grpc.UnaryServerInterceptor {
	return grpcRecovery.UnaryServerInterceptor(
		grpcRecovery.WithRecoveryHandlerContext(newGRPCRecoveryHandler(logger, enableStack)),
	)
}

func GRPCRecoverStreamServerInterceptor(logger *zerolog.Logger, enableStack bool) grpc.StreamServerInterceptor {
	return grpcRecovery.StreamServerInterceptor(
		grpcRecovery.WithRecoveryHandlerContext(newGRPCRecoveryHandler(logger, enableStack)),
	)
}

func newGRPCRecoveryHandler(logger *zerolog.Logger, enableStack bool) grpcRecovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p interface{}) error {
		err := errors.Errorf("%v", p)
		span := trace.SpanContextFromContext(ctx)

		logger.
			WithLevel(zerolog.PanicLevel).
			Stack().
			Err(err).
			Str(TypeKey, "panic").
			Str(TraceIDKey, span.TraceID().String()).
			Str(TraceFlagsKey, span.TraceFlags().String()).
			Str(ClientAddressKey, grpcRemoteIP(ctx)).
			Send()

		if enableStack {
			return status.Errorf(codes.Internal, "%+v", errors.WithStack(err))
		}

		return status.Errorf(codes.Internal, "Internal Server Error")
	}
}
