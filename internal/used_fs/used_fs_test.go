package usedfs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Av4JhG/system-monitoring/internal/sm"
)

func TestUsedFS(t *testing.T) {
	ctx := context.Background()

	_, err := Collect(ctx, sm.StartMetric)
	require.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	data, err := Collect(ctx, sm.GetMetric)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Greater(t, len(data), 0)
	for _, fs := range data {
		require.GreaterOrEqual(t, fs.UsedSpace, 0.0)
		require.GreaterOrEqual(t, fs.UsedInode, 0.0)
	}

	_, err = Collect(ctx, sm.StopMetric)
	require.NoError(t, err)
}

func TestUsedFSStartWithCanceledContext(t *testing.T) {
	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = Collect(startCtx, sm.StartMetric)

	ctx := context.Background()
	_, err := Collect(ctx, sm.StopMetric)
	require.NoError(t, err)
}
