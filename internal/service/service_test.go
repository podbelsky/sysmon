package service_test

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/podbelsky/sysmon/internal/config"
	"github.com/podbelsky/sysmon/internal/model"
	"github.com/podbelsky/sysmon/internal/service"
	"github.com/podbelsky/sysmon/internal/stat"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const mockValueNum = 3
const delta = 0.001

func TestServiceLogic(t *testing.T) {
	confTime := config.Time{
		Snap:  time.Second,
		Clean: time.Second * 30,
		Store: time.Minute,
	}
	confStat := config.Stat{
		LA:       true,
		CPU:      false,
		DiskLoad: false,
		DiskUse:  false,
	}
	collector := &MockStat{}
	logger := zerolog.Nop()

	var strangeLock sync.Mutex // added lock due race detection

	srv, err := service.NewService(confTime, confStat, collector, &logger)
	if err != nil {
		log.Fatalf("error creating service sysmon  %v", err)
	}

	log.Println("Start with duration = 6,6 sec (expecting stat dept = 6+1)")
	ticker := time.NewTicker(time.Second*6 + time.Millisecond*600)
	signalFailedStart := make(chan struct{})

	var errStart, errStop error
	go func() {
		strangeLock.Lock()

		if errStart = srv.Start(); errStart != nil {
			log.Println("failed to start service: " + err.Error())
		}
		strangeLock.Unlock()
		signalFailedStart <- struct{}{}
	}()

	select {
	case <-signalFailedStart:
	case <-ticker.C:
		t.Run("Average time more then having data", func(t *testing.T) {
			snapshot := srv.GetAverageStat(20)
			require.Nil(t, snapshot)
		})

		t.Run("Average 2 sec", func(t *testing.T) {
			data := srv.GetAverageStat(2)

			require.NotNil(t, data)

			require.InDelta(t, 6.000, data[stat.LA].Data[0].Dec, delta)
			require.InDelta(t, 12.000, data[stat.LA].Data[1].Dec, delta)
			require.InDelta(t, 18.000, data[stat.LA].Data[2].Dec, delta)
		})

		t.Run("Average 6 sec", func(t *testing.T) {
			data := srv.GetAverageStat(6)

			require.NotNil(t, data)

			require.InDelta(t, 5.666, data[stat.LA].Data[0].Dec, delta)
			require.InDelta(t, 11.333, data[stat.LA].Data[1].Dec, delta)
			require.InDelta(t, 17.000, data[stat.LA].Data[2].Dec, delta)
		})

		if errStop = srv.Stop(); errStop != nil {
			log.Println("failed to stop service: " + err.Error())
		}
	}
	log.Printf("Exiting sysmon service\n")

	strangeLock.Lock()
	require.NoError(t, errStart)
	require.NoError(t, errStop)
	strangeLock.Unlock()
}

type MockStat struct {
	counter int // odd and even values we are going to send different values
}

func (s *MockStat) LoadAvg() (*model.Bucket, error) {
	data := make([]model.Value, mockValueNum)
	// counter even la1=4 la2=8 la3=12
	// counter odd la1=8 la2=16 la3=24
	for i := 0; i < mockValueNum; i++ {
		data[i] = model.Value{
			Dec: float64(i+1) * 4 * float64(s.counter%2+1),
		}
	}
	s.counter++

	return &model.Bucket{
		Name: stat.LA,
		Data: data,
	}, nil
}

func (s *MockStat) CPUAvgStats() (*model.Bucket, error) {
	return nil, nil
}

func (s *MockStat) DisksLoad() (*model.Bucket, error) {
	return nil, nil
}

func (s *MockStat) DisksUsage() (*model.Bucket, error) {
	return nil, nil
}
