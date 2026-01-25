package sm

import "context"

// LoadAvg - функция возвращающая среднюю загрузку системы.
type LoadAvg func(ctx context.Context) (*LoadAvgData, error)

// LoadAvgInterface - обертка функции LoadAvg для моков.
type LoadAvgInterface interface {
	Execute(ctx context.Context) (*LoadAvgData, error)
}

type LoadAvgAdapter struct {
	fn LoadAvg
}

func (a *LoadAvgAdapter) Execute(ctx context.Context) (*LoadAvgData, error) {
	return a.fn(ctx)
}

func NewLoadAvg(fn LoadAvg) LoadAvgInterface {
	return &LoadAvgAdapter{fn: fn}
}

// LoadAvgData содержит метрики средней загрузки системы.
type LoadAvgData struct {
	Load1  float64
	Load5  float64
	Load15 float64
}
