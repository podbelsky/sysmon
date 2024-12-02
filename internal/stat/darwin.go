//go:build darwin

package stat

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/podbelsky/sysmon/internal/model"
)

func loadAvg() (*model.Bucket, error) {
	out, err := exec.Command(`sysctl`, `-n`, `vm.loadavg`).Output()
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(strings.Replace(string(out), ",", ".", -1))
	if len(fields) < 4 {
		return nil, errors.New("error parsing loadavg")
	}

	loadAvg1, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}

	data := make([]model.Value, 3)
	data[0] = model.Value{Dec: loadAvg1}

	loadAvg5, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}

	data[1] = model.Value{Dec: loadAvg5}

	loadAvg15, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}

	data[2] = model.Value{Dec: loadAvg15}

	return &model.Bucket{
		Name: LA,
		Data: data,
	}, nil
}

func cpuAvgStats() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}

func disksLoad() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}

func disksUsage() (*model.Bucket, error) {
	return nil, ErrNotImplemented
}
