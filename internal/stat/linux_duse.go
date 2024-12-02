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

var ErrorDiffValueDiskLoad = errors.New("diff disks set or result empty")

type DiskUsage struct {
	FileSystem string
	Type       string
	Total      uint64
	Used       uint64
	Available  uint64
	UsedPer    uint64
	MountedOn  string
}

func disksUsage() (*model.Bucket, error) {
	disksMapInode, orderInode, err := disksUsageQuery("-iTP")
	if err != nil {
		return nil, err
	}

	disksMapMb, orderMb, err := disksUsageQuery("-mTP")
	if err != nil {
		return nil, err
	}

	if len(orderInode) != len(orderMb) || len(orderInode) == 0 {
		return nil, ErrorDiffValueDiskLoad
	}

	data := make([]model.Value, 0, len(disksMapInode))

	for i, v := range orderInode {
		if orderInode[i] != orderMb[i] {
			return nil, ErrorDiffValueDiskLoad
		}

		diskMb, ok := disksMapMb[v]
		if !ok {
			return nil, err
		}

		diskInode, ok := disksMapInode[v]
		if !ok {
			return nil, err
		}

		// "Mount Point", "File System","Usage Inodes", "Usage Inode Percent", "Usage Mb", "Usage Mb Percent"
		data = append(data,
			model.Value{
				Str: v,
			}, model.Value{
				Str: diskInode.FileSystem,
			}, model.Value{
				Dec: float64(diskInode.Used),
			}, model.Value{
				Dec: float64(diskInode.UsedPer),
			}, model.Value{
				Dec: float64(diskMb.Used),
			}, model.Value{
				Dec: float64(diskMb.UsedPer),
			},
		)

		delete(disksMapInode, v)
		delete(disksMapMb, v)
	}

	if len(disksMapInode) != 0 || len(disksMapMb) != 0 {
		return nil, ErrorDiffValueDiskLoad
	}

	return &model.Bucket{
		Name: DUse,
		Data: data,
	}, nil
}

func disksUsageQuery(args string) (map[string]DiskUsage, []string, error) {
	diskUsageMap := make(map[string]DiskUsage)
	order := make([]string, 0, 5)

	df, err := exec.LookPath("df")
	if err != nil {
		return nil, nil, err
	}

	out, err := exec.Command(df, args).Output()
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	// Filter the header
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		diskUsage, err := disksUsageQueryParse(line)
		if err != nil {
			return nil, nil, err
		}

		diskUsageMap[diskUsage.MountedOn] = diskUsage
		order = append(order, diskUsage.MountedOn)
	}

	return diskUsageMap, order, nil
}

func disksUsageQueryParse(usage string) (diskUsage DiskUsage, err error) {
	diskUsage = DiskUsage{}

	fields := strings.Fields(usage)

	// Check there are 7 fields
	if len(fields) != 7 {
		return DiskUsage{}, errors.New("couldn't parse disk usage because there aren't 7 fields")
	}

	// Parse fields
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		switch i {
		case 0:
			diskUsage.FileSystem = field
		case 1:
			diskUsage.Type = field
		case 2:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Total = value
		case 3:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Used = value
		case 4:
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return DiskUsage{}, err
			}
			diskUsage.Available = value
		case 5:
			// Trim trailing '%'
			if last := len(field) - 1; last >= 0 && field[last] == '%' {
				field = field[:last]
			}
			var value uint64
			if field == "-" {
				value = 0
			} else {
				value, err = strconv.ParseUint(field, 10, 64)
				if err != nil {
					return DiskUsage{}, err
				}
			}
			diskUsage.UsedPer = value
		case 6:
			diskUsage.MountedOn = field
		}
	}

	return diskUsage, nil
}
