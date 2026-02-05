package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	conf "github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/cpu"
	loadavg "github.com/Av4JhG/system-monitoring/internal/load_avg"
	loaddisks "github.com/Av4JhG/system-monitoring/internal/load_disks"
	"github.com/Av4JhG/system-monitoring/internal/mocks"
	"github.com/Av4JhG/system-monitoring/internal/sm"
	usedfs "github.com/Av4JhG/system-monitoring/internal/used_fs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestCollectorStartStop(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := conf.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.MockLogger)
	log.On("Debug", "collector is stopped")

	collectors := sm.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toClientsCh := make(sm.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toClientsCh, 0)
}

func TestCollectorStartWithCanceledContext(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := conf.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)

	collectors := sm.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toClientsCh := make(sm.CollectorToClientsCh, 1)

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toClientsCh, 0)
}

func TestCollectorTick(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := conf.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything)
	log.On("Debug", "tick ", mock.Anything)

	laData := &sm.LoadAvgData{
		Load1:  1,
		Load5:  2,
		Load15: 3,
	}
	LoadAvg := new(mocks.MockLoadAvgInterface)
	LoadAvg.On("Execute", mock.Anything).Return(laData, nil)

	cpuData := &sm.CPUData{
		User:   0.1,
		System: 0.2,
		Idle:   0.3,
	}
	CPU := new(mocks.MockCPUInterface)
	CPU.On("Execute", mock.Anything, mock.Anything).Return(cpuData, nil)

	ldData := sm.LoadDisksData{
		{
			Name:    "sda",
			Tps:     7,
			KBRead:  8,
			KBWrite: 9,
		},
	}
	LoadDisks := new(mocks.MockLoadDisksInterface)
	LoadDisks.On("Execute", mock.Anything, mock.Anything).Return(ldData, nil)

	fsData := sm.UsedFSData{
		{
			Path:      "/",
			UsedSpace: 13,
			UsedInode: 31,
		},
	}
	UsedFS := new(mocks.MockUsedFSInterface)
	UsedFS.On("Execute", mock.Anything, mock.Anything).Return(fsData, nil)

	collectors := sm.MetricCollectors{
		LoadAvg:   LoadAvg.Execute,
		CPU:       CPU.Execute,
		LoadDisks: LoadDisks.Execute,
		UsedFS:    UsedFS.Execute,
	}

	toClientsCh := make(sm.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	time.Sleep(50 * time.Millisecond)

	time.Sleep(time.Second)
	data := <-toClientsCh
	// текущая секунда еще не заполнена, предыдущих нет - статистика должна быть пустой
	require.Len(t, data.Points, 0)

	time.Sleep(time.Second)
	data = <-toClientsCh
	require.Len(t, data.Points, 1)
	for _, point := range data.Points {
		require.Equal(t, laData, point.LoadAvg)
		require.Equal(t, cpuData, point.CPU)
		require.Equal(t, ldData, point.LoadDisks)
		require.Equal(t, fsData, point.UsedFS)
	}

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
}

func TestCollectorTickWithErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := conf.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)
	log.On("Debug", "tick ", mock.Anything)

	laErr := errors.New("LoadAvg Error")
	LoadAvg := new(mocks.MockLoadAvgInterface)
	LoadAvg.On("Execute", mock.Anything).Return(nil, laErr)

	cpuErr := errors.New("CPU Error")
	CPU := new(mocks.MockCPUInterface)
	CPU.On("Execute", mock.Anything, sm.StartMetric).Return(nil, nil)
	CPU.On("Execute", mock.Anything, sm.StopMetric).Return(nil, nil)
	CPU.On("Execute", mock.Anything, sm.GetMetric).Return(nil, cpuErr)

	ldErr := errors.New("LoadDisks Error")
	LoadDisks := new(mocks.MockLoadDisksInterface)
	LoadDisks.On("Execute", mock.Anything, sm.StartMetric).Return(nil, nil)
	LoadDisks.On("Execute", mock.Anything, sm.StopMetric).Return(nil, nil)
	LoadDisks.On("Execute", mock.Anything, sm.GetMetric).Return(nil, ldErr)

	fsErr := errors.New("UsedFS Error")
	UsedFS := new(mocks.MockUsedFSInterface)
	UsedFS.On("Execute", mock.Anything, sm.StartMetric).Return(nil, nil)
	UsedFS.On("Execute", mock.Anything, sm.StopMetric).Return(nil, nil)
	UsedFS.On("Execute", mock.Anything, sm.GetMetric).Return(nil, fsErr)

	collectors := sm.MetricCollectors{
		LoadAvg:   LoadAvg.Execute,
		CPU:       CPU.Execute,
		LoadDisks: LoadDisks.Execute,
		UsedFS:    UsedFS.Execute,
	}

	toClientsCh := make(sm.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	time.Sleep(50 * time.Millisecond)

	time.Sleep(time.Second)
	data := <-toClientsCh
	// текущая секунда еще не заполнена, предыдущих нет - статистика должна быть пустой
	require.Len(t, data.Points, 0)

	// все коллекторы вернули ошибки, данные должны быть пустые
	time.Sleep(time.Second)
	data = <-toClientsCh
	require.Len(t, data.Points, 1)
	for _, point := range data.Points {
		require.Nil(t, point.LoadAvg)
		require.Nil(t, point.CPU)
		require.Nil(t, point.LoadDisks)
		require.Nil(t, point.UsedFS)
	}

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
}
