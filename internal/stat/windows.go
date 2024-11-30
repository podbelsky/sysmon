//go:build windows

package stat

import (
	"github.com/podbelsky/sysmon/internal/model"
)

func loadAvg() (*model.Bucket, error) { return nil, ErrNotImplemented }

func cpuAvgStats() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}

func disksLoad() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}

func disksUsage() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}
