package cmd

import (
	"context"
	"log"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/build"
	"github.com/podbelsky/sysmon/internal/client"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/podbelsky/sysmon/internal/model"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func clientCmd(ctx context.Context, conf config.Config) *cobra.Command {
	flagN := "N"
	flagM := "M"

	cmd := &cobra.Command{
		Use:   "client",
		Short: "grpc client consume",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			average, _ := cmd.Flags().GetInt(flagM)
			query, _ := cmd.Flags().GetInt(flagN)
			zerolog.Ctx(ctx).Info().
				Int("period", query).
				Int("average", average).
				Msg("start " + cmd.Short)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			zerolog.Ctx(ctx).Info().Msg("stop " + cmd.Short)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := build.New(ctx, conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			average, _ := cmd.Flags().GetInt(flagM)
			query, _ := cmd.Flags().GetInt(flagN)

			cl := client.NewClient(conf.GRPCAddr(), query, average)

			err := cl.Connect()
			if err != nil {
				return errors.Wrap(err, "client connect failed")
			}

			defer func() {
				err = cl.Close()
				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("client close failed")
				}
			}()

			// Start getting data from GRPC server
			go func() {
				err = cl.GetData(func(data model.Snapshot) {
					for _, v := range data {
						log.Printf("%s %v\n", v.Name, v.Data)
					}
				})

				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("get data failed")
					cancel()
				}
			}()

			<-ctx.Done()

			return nil
		},
	}

	cmd.PersistentFlags().Int(flagN, 2, "query period")
	cmd.PersistentFlags().Int(flagM, 10, "average interval")

	return cmd
}
