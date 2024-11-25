package build

import (
	"context"

	"github.com/podbelsky/sysmon/internal/config"
)

type Builder struct {
	config   config.Config
	shutdown shutdown
}

func New(ctx context.Context, conf config.Config) *Builder {
	b := Builder{config: conf} //nolint:exhaustruct

	return &b
}
