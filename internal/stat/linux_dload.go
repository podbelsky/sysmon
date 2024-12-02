//go:build linux

package stat

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/podbelsky/sysmon/internal/model"
)

var errorSkipLine = errors.New("skip line")

type DiskLoad struct {
	Device string
	Tps    float64
	KbRW   float64
}

func disksLoad() (*model.Bucket, error) {
	dl, err := disksLoadQuery()
	if err != nil {
		return nil, err
	}

	data := make([]model.Value, 0, len(dl))
	for _, v := range dl {
		// "(Device)", "(Tps)", "(Kbps)"
		data = append(data, model.Value{
			Str: v.Device,
		}, model.Value{
			Dec: v.Tps,
		}, model.Value{
			Dec: v.KbRW,
		})
	}

	return &model.Bucket{
		Name: DLoad,
		Data: data,
	}, nil
}

func disksLoadQuery() ([]DiskLoad, error) {
	iostat, err := exec.LookPath("iostat")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(iostat, "-yd", "1", "1").Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	disksload := make([]DiskLoad, 0, 2)
	for scanner.Scan() {
		line := scanner.Text()
		diskload, err := disksLoadQueryParse(line)
		if err != nil {
			if errors.Is(err, errorSkipLine) {
				continue
			}
			return nil, err
		}
		disksload = append(disksload, diskload)
	}
	return disksload, nil
}

func disksLoadQueryParse(load string) (DiskLoad, error) {
	fields := strings.Fields(load)
	if len(fields) == 0 {
		return DiskLoad{}, errorSkipLine
	}
	if strings.Contains(fields[0], "loop") {
		return DiskLoad{}, errorSkipLine
	}
	var diskLoad DiskLoad
	diskLoad.Device = fields[0]
	tps, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	diskLoad.Tps = tps
	kbr, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	kbw, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	diskLoad.KbRW = kbr + kbw

	return diskLoad, nil
}
