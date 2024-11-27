package external

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/podbelsky/sysmon/pkg/monitor/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/peer"
)

type Service interface {
	Start() error
	Stop() error
	GetAverageStat(ctx context.Context, average int) interface{}
}

type Server struct {
	service Service
	logger  *zerolog.Logger
	v1.UnimplementedMonitorAPIServer
}

func NewServer(s Service, l *zerolog.Logger) *Server {
	return &Server{ //nolint:exhaustruct
		service: s,
		logger:  l,
	}
}

func (s *Server) GetStat(request *v1.GetStatRequest, server v1.MonitorAPI_GetStatServer) error {
	ticker := time.NewTicker(time.Second * time.Duration(request.GetPeriod()))
	defer func() {
		ticker.Stop()
	}()

	var address string
	if p, ok := peer.FromContext(server.Context()); ok {
		address = p.Addr.String()
	}

	s.logger.Info().
		Int32("average", request.GetAverage()).
		Int32("period", request.GetPeriod()).
		Str("address", address).
		Msg("New gRPC connection")

	defer s.logger.Info().Str("address", address).Msg("End gRPC connection")

	send := func() error {
		r := s.service.GetAverageStat(server.Context(), int(request.GetAverage()))
		b, err := json.Marshal(&r)
		if err != nil {
			return err
		}

		return server.Send(&v1.GetStatResponse{Data: b})
	}

	if err := send(); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := send(); err != nil {
				return err
			}
		case <-server.Context().Done():
			return nil
		}
	}
}
