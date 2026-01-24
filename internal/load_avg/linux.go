//go:build linux
// +build linux

package loadavg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Av4JhG/system-monitoring/internal/core"
	"github.com/Av4JhG/system-monitoring/internal/sm"
)

// Collect позволяет управлять получением метрик средней загрузки системы.
func Collect(_ context.Context) (*sm.LoadAvgData, error) {
	content, err := core.ReadProcFile("loadavg")
	if err != nil {
		return nil, fmt.Errorf("cannot read the loadavg file: %w", err)
	}

	return parseLoadAvg(content)
}

func parseLoadAvg(content []string) (*sm.LoadAvgData, error) {
	values := strings.Fields(content[0])

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load1 field: %w", err)
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load5 field: %w", err)
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load15 field: %w", err)
	}

	return &sm.LoadAvgData{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}
