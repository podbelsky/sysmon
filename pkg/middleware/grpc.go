package middleware

import (
	"context"
	"encoding/binary"
	"path"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func grpcRemoteIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return ""
}

func grpcMethodName(method string) string {
	return path.Base(method)
}

func grpcUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("user-agent"), "")
	}

	return ""
}

func grpcStatusCode(err error) int {
	return int(status.Convert(err).Code())
}

func grpcMessageBytesCount(message interface{}) int {
	if pb, ok := message.(proto.Message); ok {
		if b, err := protojson.Marshal(pb); err == nil {
			return binary.Size(b)
		}
	}

	return 0
}
