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

func TestLoadDisks(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)

	ldData := sm.LoadDisksData{
		{
			Name:    "sda",
			Tps:     5,
			KBRead:  7,
			KBWrite: 12,
		},
	}

	collector := func(_ context.Context, _ sm.MetricCommand) (sm.LoadDisksData, error) {
		return ldData, nil
	}

	loadDisksCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, ldData, point.LoadDisks)
}

func TestLoadDisksError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.MockLogger)
	log.On("Debug", mock.Anything, mock.Anything)

	collector := func(_ context.Context, _ sm.MetricCommand) (sm.LoadDisksData, error) {
		return nil, fmt.Errorf("cannot parse iostat line")
	}

	loadDisksCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadDisks)
}
