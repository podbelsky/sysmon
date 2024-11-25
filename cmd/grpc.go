package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/build"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/spf13/cobra"
	grpclib "google.golang.org/grpc"
)

func grpcCmd(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "start grpc server listening",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := build.New(ctx, conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			listener, err := builder.Listener()
			if err != nil {
				return errors.Wrap(err, "start network listener")
			}

			server, err := builder.MonitorGRPCServer(ctx)
			if err != nil {
				return errors.Wrap(err, "build grpc server")
			}

			if err = server.Serve(listener); !errors.Is(err, grpclib.ErrServerStopped) {
				return errors.Wrap(err, "run grpc server")
			}

			<-ctx.Done()

			return nil
		},
	}
}
