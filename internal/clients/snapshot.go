package clients

import (
	"time"

	"github.com/Av4JhG/system-monitoring/internal/sm"
)

func makeSnapshot(data *sm.MetricsData, m int) *sm.Stats {
	result := &sm.Stats{
		Time: data.Time,
	}

	from := data.Time.Add(time.Duration(-m) * time.Second)
	points := make([]*sm.Point, 0, len(data.Points))

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}
		points = append(points, point)
	}

	fillCPU(result, points)

	return result
}

func fillCPU(result *sm.Stats, points []*sm.Point) {
	countCPU := 0
	cpuUser := 0.0
	cpuSystem := 0.0
	cpiIdle := 0.0

	for _, point := range points {
		if point.CPU != nil {
			countCPU++
			cpuUser += point.CPU.User
			cpuSystem += point.CPU.System
			cpiIdle += point.CPU.Idle
		}
	}

	if countCPU > 0 {
		result.CPU = &sm.CPUData{
			User:   cpuUser / float64(countCPU),
			System: cpuSystem / float64(countCPU),
			Idle:   cpiIdle / float64(countCPU),
		}
	}
}
