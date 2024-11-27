package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"

	"github.com/podbelsky/sysmon/internal/external"
	v1 "github.com/podbelsky/sysmon/pkg/monitor/v1"
)

func (b *Builder) MonitorGRPCServer(ctx context.Context) (*grpc.Server, external.Service, error) {
	service, err := b.Service(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "build monitor service")
	}

	grpcServer := b.gRPCServer(ctx)

	apiServer := external.NewServer(service, b.logger)

	v1.RegisterMonitorAPIServer(grpcServer, apiServer)

	return grpcServer, service, nil
}
