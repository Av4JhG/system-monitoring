//go:build linux
// +build linux

package cpu

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Av4JhG/system-monitoring/internal/core"
)

func TestParseCPU(t *testing.T) {
	content, err := os.ReadFile("./test_data/stat_1")
	require.NoError(t, err)

	data, err := parseCPU(core.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, float64(631341+1284+109025+3744304+11237+1685), data.total)
	require.Equal(t, 631341.0, data.user)
	require.Equal(t, 109025.0, data.system)
	require.Equal(t, 3744304.0, data.idle)
}

func TestParseCPUFail(t *testing.T) {
	content, err := os.ReadFile("./test_data/stat_2")
	require.NoError(t, err)

	_, err = parseCPU(core.SplitLines(string(content)))
	require.ErrorIs(t, err, strconv.ErrSyntax)
}
