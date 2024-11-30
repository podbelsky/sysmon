package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/external"
	"github.com/podbelsky/sysmon/internal/service"
	"github.com/podbelsky/sysmon/internal/stat"
)

func (b *Builder) Service(_ context.Context) (external.Service, error) {
	collector := &stat.System{}
	srv, err := service.NewService(b.config.Time, b.config.Stat, collector, b.logger)
	if err != nil {
		return nil, err
	}

	b.shutdown.add(func(_ context.Context) error {
		if err = srv.Stop(); err != nil {
			return errors.Wrap(err, "stop monitor service")
		}

		return nil
	})

	return srv, nil
}
