package build

import (
	"context"

	"google.golang.org/grpc"

	v1 "github.com/podbelsky/sysmon/pkg/monitor/v1"
)

func (b *Builder) MonitorGRPCServer(ctx context.Context) (*grpc.Server, error) {
	//service, err := b.Service(ctx)
	//if err != nil {
	//	return nil, errors.Wrap(err, "build service")
	//}

	s := b.gRPCServer(ctx)

	// TODO implement
	apiServer := &v1.UnimplementedMonitorAPIServer{}

	v1.RegisterMonitorAPIServer(s, apiServer)

	return s, nil
}
