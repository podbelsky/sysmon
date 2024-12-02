//go:build linux

package stat

import (
	"os"
	"strconv"
	"strings"

	"github.com/podbelsky/sysmon/internal/model"
)

func loadAvg() (*model.Bucket, error) {
	file, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(file))
	loadAvg1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, err
	}

	data := make([]model.Value, 3)
	data[0] = model.Value{Dec: loadAvg1}

	loadAvg5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}

	data[1] = model.Value{Dec: loadAvg5}

	loadAvg15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}

	data[0] = model.Value{Dec: loadAvg15}

	return &model.Bucket{
		Name: LA,
		Data: data,
	}, nil
}
