package sm

import "context"

// CPU - функция возвращающая среднюю загрузку cpu.
type CPU func(ctx context.Context, action MetricCommand) (*CPUData, error)

// CPUInterface - обертка функции CPU для моков.
type CPUInterface interface {
	Execute(ctx context.Context, action MetricCommand) (*CPUData, error)
}

type CPUAdapter struct {
	fn CPU
}

func (a *CPUAdapter) Execute(ctx context.Context, action MetricCommand) (*CPUData, error) {
	return a.fn(ctx, action)
}

func NewCPU(fn CPU) CPUInterface {
	return &CPUAdapter{fn: fn}
}

// CPUData содержит метрики средней загрузки cpu. В процентах.
type CPUData struct {
	User   float64
	System float64
	Idle   float64
}
