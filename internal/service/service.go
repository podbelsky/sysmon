package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/podbelsky/sysmon/internal/model"
	"github.com/podbelsky/sysmon/internal/stat"
	"github.com/rs/zerolog"
)

var (
	ErrStoreLessSnap = errors.New("store interval cannot be less then snap interval")
)

type Service struct {
	cfg config.Time

	data    model.Storage
	workers []stat.Fn

	done chan struct{}
	wg   sync.WaitGroup

	logger *zerolog.Logger
}

func NewService(t config.Time, s config.Stat, collector stat.Collector, l *zerolog.Logger) (*Service, error) {
	if t.Store < t.Snap {
		return nil, ErrStoreLessSnap
	}

	workers := make([]stat.Fn, 0)
	if s.LA {
		workers = append(workers, collector.LoadAvg)
	}
	if s.CPU {
		workers = append(workers, collector.CPUAvgStats)
	}
	if s.DiskLoad {
		workers = append(workers, collector.DisksLoad)
	}
	if s.DiskUse {
		workers = append(workers, collector.DisksUsage)
	}

	maxSize := int(t.Store/t.Snap) + 1

	return &Service{
		cfg:     t,
		workers: workers,
		data: model.Storage{
			History: make(map[int]model.Snapshot, maxSize),
			Limit:   maxSize,
		},
		done: make(chan struct{}),
		wg:   sync.WaitGroup{},

		logger: l,
	}, nil
}

func (s *Service) Start() error {
	s.logger.Info().
		Str("time.snap", s.cfg.Snap.String()).
		Str("time.clean", s.cfg.Clean.String()).
		Str("time.store", s.cfg.Store.String()).
		Msg("Starting monitor service...")

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

	close(s.done)
	s.wg.Wait()

	return nil
}

func (s *Service) GetAverageStat(average int) model.Snapshot {
	if average <= 0 {
		s.logger.Error().Int("average", average).Msg("average must be gt zero")

		return nil
	}

	s.data.Lock.RLock()
	defer s.data.Lock.RUnlock()

	size := len(s.data.Index)
	snapPeriodSecond := int(s.cfg.Snap / time.Second)
	if snapPeriodSecond*size < average {
		s.logger.Warn().Msg("stat dept lower than average period")

		return nil
	}

	elementsNum := average / snapPeriodSecond
	if elementsNum == 0 {
		elementsNum = 1
	}

	res := make(model.Snapshot, len(s.workers))

	for i := 1; i <= elementsNum; i++ {
		offset := elementsNum - i
		snapshot, ok := s.data.History[s.data.Index[offset]]
		if !ok {
			continue
		}

		for name, bucket := range snapshot {
			sumBucket, ok := res[name]

			if !ok {
				res[name] = bucket
				continue
			}

			for j, val := range bucket.Data {
				sumBucket.Data[j].Dec += val.Dec
			}
		}
	}

	for name, bucket := range res {
		for j, val := range bucket.Data {
			bucket.Data[j].Dec = val.Dec / float64(elementsNum)
		}

		bucket.Name = fmt.Sprintf("%s(%d)", bucket.Name, elementsNum)
		res[name] = bucket
	}

	return res
}

func (s *Service) snap() {
	ticker := time.NewTicker(s.cfg.Snap)
	defer func() {
		ticker.Stop()
	}()

	collect := func() {
		snapshot := make(model.Snapshot, len(s.workers))

		wg := &sync.WaitGroup{}
		wg.Add(len(s.workers))

		for _, fn := range s.workers {
			go func() {
				defer wg.Done()

				bucket, err := fn()
				if err != nil {
					if !errors.Is(err, stat.ErrNotImplemented) {
						s.logger.Err(err).Msg("get stat failed")
					}

					return
				}

				snapshot[bucket.Name] = *bucket

				s.logger.Trace().Msgf("%v", bucket)
			}()
		}

		wg.Wait()

		s.data.Lock.Lock()

		s.data.Counter++
		s.data.Index = append(s.data.Index, s.data.Counter)
		s.data.History[s.data.Counter] = snapshot

		s.data.Lock.Unlock()

		s.logger.Debug().
			Int("stat.size", len(s.data.History)).
			//Str("snapshot", fmt.Sprintf("%v", snapshot)).
			Msg("snap")
	}
	collect()

	for {
		select {
		case <-ticker.C:
			collect()
		case <-s.done:
			return
		}
	}
}

func (s *Service) clean() {
	ticker := time.NewTicker(s.cfg.Clean)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			s.data.Lock.Lock()

			num := len(s.data.Index) - s.data.Limit
			if num > 0 {
				for i := 0; i < num; i++ {
					idx := s.data.Index[0]
					delete(s.data.History, idx)
					s.data.Index = s.data.Index[1:len(s.data.Index)]
				}
			}

			s.data.Lock.Unlock()
			s.logger.Debug().
				Int("data.size", len(s.data.History)).
				Msg("clean")

		case <-s.done:
			return
		}
	}
}
