package sm

import "context"

// UsedFS - функция возвращающая использование файловых систем.
type UsedFS func(ctx context.Context, action MetricCommand) (UsedFSData, error)

// LoadDisksInterface - обертка функции LoadDisks для моков.
type UsedFSInterface interface {
	Execute(ctx context.Context, action MetricCommand) (UsedFSData, error)
}

type UsedFSAdapter struct {
	fn UsedFS
}

func (a *UsedFSAdapter) Execute(ctx context.Context, action MetricCommand) (UsedFSData, error) {
	return a.fn(ctx, action)
}

func NewUsedFS(fn UsedFS) UsedFSInterface {
	return &UsedFSAdapter{fn: fn}
}

// UsedFSData - слайс информации об использовании файловых систем.
type UsedFSData []FSData

// FSData содержит информацию об использовании файловой системы.
type FSData struct {
	Path      string
	UsedSpace float64
	UsedInode float64
}
