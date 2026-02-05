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

func TestUsedFS(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)

	ufData := sm.UsedFSData{
		{
			Path:      "/",
			UsedSpace: 12.3,
			UsedInode: 7.77,
		},
	}

	collector := func(_ context.Context, _ sm.MetricCommand) (sm.UsedFSData, error) {
		return ufData, nil
	}

	usedFSCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, ufData, point.UsedFS)
}

func TestUsedFSError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)

	collector := func(_ context.Context, _ sm.MetricCommand) (sm.UsedFSData, error) {
		return nil, fmt.Errorf("cannot parse df line")
	}

	usedFSCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.UsedFS)
}
