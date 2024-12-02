package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/podbelsky/sysmon/internal/version"
	"github.com/spf13/cobra"
)

func Run(ctx context.Context, conf config.Config) error {
	root := &cobra.Command{ //nolint:exhaustruct
		Version: version.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	root.AddCommand(
		grpcCmd(ctx, conf),
		clientCmd(ctx, conf),
	)

	return errors.Wrap(root.ExecuteContext(ctx), "run application")
}
