package middleware

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func GRPCAccessLogUnaryServerInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		logErr := ctx.Err()
		if err != nil {
			logErr = err
		}

		grpcAccessLogEvent(ctx, logger, info.FullMethod, start, logErr).
			Int(RPCServerRequestSizeKey, grpcMessageBytesCount(req)).
			Int(RPCServerResponseSizeKey, grpcMessageBytesCount(resp)).
			Send()

		return resp, err
	}
}

func GRPCAccessLogStreamServerInterceptor(logger *zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		start := time.Now()
		err := handler(srv, ss)

		logErr := ctx.Err()
		if err != nil {
			logErr = err
		}

		grpcAccessLogEvent(ctx, logger, info.FullMethod, start, logErr).Send()

		return err
	}
}

func grpcAccessLogEvent(
	ctx context.Context,
	logger *zerolog.Logger,
	method string,
	start time.Time,
	err error,
) *zerolog.Event {
	//nolint:zerologlint
	event := logger.WithLevel(zerolog.NoLevel).
		Err(err).
		Str(TypeKey, "access").
		Int64(RPCServerDurationKey, time.Since(start).Milliseconds()).
		Str(UserAgentOriginalKey, grpcUserAgent(ctx)).
		Str(RPCMethodKey, grpcMethodName(method)).
		Int(RPCGRPCStatusCodeKey, grpcStatusCode(err)).
		Str(ClientAddressKey, grpcRemoteIP(ctx))

	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		event = event.
			Str(TraceIDKey, span.TraceID().String()).
			Str(TraceFlagsKey, span.TraceFlags().String())
	}

	return event
}
