package service

import (
	"context"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/podbelsky/sysmon/internal/model"
	"github.com/rs/zerolog"
)

var (
	ErrStoreLessSnap = errors.New("store interval cannot be less then snap interval")
)

type Service struct {
	cfgTime config.Time
	cfgStat config.Stat

	exitChan chan struct{}
	wg       sync.WaitGroup

	logger *zerolog.Logger
}

func NewService(t config.Time, s config.Stat, l *zerolog.Logger) (*Service, error) {
	if t.Store < t.Snap {
		return nil, ErrStoreLessSnap
	}

	return &Service{
		cfgTime: t,
		cfgStat: s,

		exitChan: make(chan struct{}),
		wg:       sync.WaitGroup{},

		logger: l,
	}, nil
}

func (s *Service) Start() error {
	s.logger.Info().Msg("Starting monitor service...")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.snap()
	}()

	go func() {
		defer wg.Done()
		s.clean()
	}()

	wg.Wait()

	return nil
}

func (s *Service) Stop() error {
	s.logger.Info().Msg("Stopping monitor service...")

	close(s.exitChan)
	//s.stat.Close()
	s.wg.Wait()

	return nil
}

func (s *Service) GetAverageStat(ctx context.Context, average int) interface{} {
	// TODO implement

	r := []model.DataToClientStamp{
		{
			Name:      "test",
			IdxHeader: 0,
			Data: [][]string{
				{"123", "456", "789"},
				{"222", "333", "444"},
			},
		},
	}

	return r
}

func (s *Service) snap() {
	ticker := time.NewTicker(s.cfgTime.Snap)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			s.logger.Debug().Msg("snap")
			// TODO implement

		case <-s.exitChan:
			return
		}
	}
}

func (s *Service) clean() {
	ticker := time.NewTicker(s.cfgTime.Clean)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			s.logger.Debug().Msg("clean")
			// TODO implement

			//s.data.Lock.Lock()
			//num := len(s.data.Index) - s.data.MaxElements
			//if num > 0 {
			//	for i := 0; i < num; i++ {
			//		idx := s.data.Index[0]
			//		delete(s.data.Elements, idx)
			//		s.data.Index = s.data.Index[1:len(s.data.Index)]
			//	}
			//}
			//s.data.Lock.Unlock()
		case <-s.exitChan:
			return
		}
	}
}
