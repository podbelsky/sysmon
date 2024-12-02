//go:build linux

package stat

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/podbelsky/sysmon/internal/model"
)

func cpuAvgStats() (*model.Bucket, error) {

	total1, values1, err := getCPUStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 100)

	total2, values2, err := getCPUStats()
	if err != nil {
		return nil, err
	}

	delta := abs(total2 - total1)

	size := len(values1)
	data := make([]model.Value, size)

	for i := 0; i < size; i++ {
		data[i] = model.Value{
			Dec: abs(values2[i]-values1[i]) * 100 / delta,
		}
	}

	return &model.Bucket{
		Name: CPU,
		Data: data,
	}, nil
}

func getCPUStats() (total float64, stat []float64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0.0, nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())

	for i := range fields {
		if i == 0 {
			continue
		}

		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return 0.0, nil, err
		}

		total += value

		switch {
		case i == 1: // User
			stat = append(stat, value)
		case i == 3: // System
			stat = append(stat, value)
		case i == 4: // Idle
			stat = append(stat, value)
		}
	}

	return total, stat, nil
}

func abs(f float64) float64 {
	if f > 0 {
		return f
	}
	return -1 * f
}
