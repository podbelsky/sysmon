package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/external"
	"github.com/podbelsky/sysmon/internal/service"
)

func (b *Builder) Service(ctx context.Context) (external.Service, error) {
	srv, err := service.NewService(b.config.Time, b.config.Stat, b.logger)
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
