package collector

import (
	"context"
	"fmt"
	"testing"

	"github.com/Av4JhG/system-monitoring/internal/mocks"
	"github.com/Av4JhG/system-monitoring/internal/sm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoadAvg(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)

	loadAvg := sm.LoadAvgData{
		Load1:  0.1,
		Load5:  0.2,
		Load15: 0.3,
	}

	collector := func(_ context.Context) (*sm.LoadAvgData, error) {
		return &loadAvg, nil
	}

	loadavgCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, &loadAvg, point.LoadAvg)
}

func TestLoadAvgError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)

	collector := func(_ context.Context) (*sm.LoadAvgData, error) {
		return nil, fmt.Errorf("cannot read the loadavg file")
	}

	loadavgCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadAvg)
}
