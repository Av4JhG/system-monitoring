//go:build linux
// +build linux

package loadavg

import (
	"os"
	"strconv"
	"testing"

	"github.com/Av4JhG/system-monitoring/internal/core"
	"github.com/stretchr/testify/require"
)

func TestParseLoadAvg(t *testing.T) {
	content, err := os.ReadFile("./test_data/load_avg_1")
	require.NoError(t, err)

	data, err := parseLoadAvg(core.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, 0.60, data.Load1)
	require.Equal(t, 0.75, data.Load5)
	require.Equal(t, 0.74, data.Load15)
}

func TestParseLoadAvgFail(t *testing.T) {
	content, err := os.ReadFile("./test_data/load_avg_2")
	require.NoError(t, err)

	_, err = parseLoadAvg(core.SplitLines(string(content)))
	require.ErrorIs(t, err, strconv.ErrSyntax)
}
