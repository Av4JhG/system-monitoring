package collector

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Av4JhG/system-monitoring/internal/mocks"
	"github.com/Av4JhG/system-monitoring/internal/sm"
)

func TestCPU(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)

	cpuData := sm.CPUData{
		User:   0.1,
		System: 0.2,
		Idle:   0.3,
	}

	collector := func(_ context.Context, _ sm.MetricCommand) (*sm.CPUData, error) {
		return &cpuData, nil
	}

	cpuCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, &cpuData, point.CPU)
}

func TestCPUError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)

	collector := func(_ context.Context, _ sm.MetricCommand) (*sm.CPUData, error) {
		return nil, fmt.Errorf("cannot read the stat file")
	}

	cpuCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.CPU)
}
