package build

import (
	"context"

	"github.com/podbelsky/sysmon/internal/config"
	"github.com/rs/zerolog"
)

type Builder struct {
	config   config.Config
	shutdown shutdown
	logger   *zerolog.Logger
}

func New(ctx context.Context, conf config.Config) *Builder {
	logger := zerolog.Ctx(ctx)
	b := Builder{config: conf, logger: logger} //nolint:exhaustruct

	return &b
}
