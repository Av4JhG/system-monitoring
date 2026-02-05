package sm

import "context"

// LoadDisks - функция возвращающая загрузку дисков.
type LoadDisks func(ctx context.Context, action MetricCommand) (LoadDisksData, error)

// LoadDisksInterface - обертка функции LoadDisks для моков.
type LoadDisksInterface interface {
	Execute(ctx context.Context, action MetricCommand) (LoadDisksData, error)
}

type LoadDisksAdapter struct {
	fn LoadDisks
}

func (a *LoadDisksAdapter) Execute(ctx context.Context, action MetricCommand) (LoadDisksData, error) {
	return a.fn(ctx, action)
}

func NewLoadDisks(fn LoadDisks) LoadDisksInterface {
	return &LoadDisksAdapter{fn: fn}
}

// LoadDisksData - слайс информации о загрузке дисков.
type LoadDisksData []DiskData

// DiskData содержит метрики загрузки дисков.
type DiskData struct {
	Name    string
	Tps     float64
	KBRead  float64
	KBWrite float64
}
