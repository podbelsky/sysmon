package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/podbelsky/sysmon/internal/model"
	v1 "github.com/podbelsky/sysmon/pkg/monitor/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Monitor struct {
	client  v1.MonitorAPIClient
	conn    *grpc.ClientConn
	addr    string
	period  int
	average int
}

func NewClient(address string, period, average int) *Monitor {
	return &Monitor{
		addr:    address,
		period:  period,
		average: average,
	}
}

func (s *Monitor) Connect() error {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:staticcheck
	if err != nil {
		return err
	}

	s.conn = conn
	s.client = v1.NewMonitorAPIClient(conn)

	return nil
}

func (s *Monitor) Close() error {
	return s.conn.Close()
}

func (s *Monitor) GetData(callback func(model.Snapshot)) error {
	DataStream, err := s.client.GetStat(context.Background(), &v1.GetStatRequest{
		Period:  int32(s.period),
		Average: int32(s.average),
	})

	if err != nil {
		return err
	}

	for {
		stream, err := DataStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		var dataModel model.Snapshot
		err = json.Unmarshal(stream.Data, &dataModel)
		if err != nil {
			return err
		}

		callback(dataModel)
	}

	return nil
}
