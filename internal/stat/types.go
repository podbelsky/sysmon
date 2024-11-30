package stat

import (
	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/model"
)

const (
	LA    = "loadAvg"
	CPU   = "cpuAvg"
	DLoad = "disksLoad"
	DUse  = "diskUsage"
)

var ErrNotImplemented = errors.New("not implemented")

type Fn func() (*model.Bucket, error)

type Collector interface {
	LoadAvg() (*model.Bucket, error)
	CPUAvgStats() (*model.Bucket, error)
	DisksLoad() (*model.Bucket, error)
	DisksUsage() (*model.Bucket, error)
}

type System struct{}

func (s *System) LoadAvg() (*model.Bucket, error) { return loadAvg() }

func (s *System) CPUAvgStats() (*model.Bucket, error) {
	return cpuAvgStats()
}

func (s *System) DisksLoad() (*model.Bucket, error) {
	return disksLoad()
}

func (s *System) DisksUsage() (*model.Bucket, error) { return disksUsage() }
