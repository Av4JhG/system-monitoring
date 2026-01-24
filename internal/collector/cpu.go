package collector

import (
	"context"
	"errors"
	"sync"

	"github.com/Av4JhG/system-monitoring/internal/sm"
)

func cpuCollect(ctx context.Context, mutex sync.Locker, ch <-chan timePoint, collector sm.CPU, log sm.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case tp := <-ch:
			func() {
				workCtx, cancel := context.WithTimeout(ctx, timeToGetMetric)
				defer cancel()

				data, err := collector(workCtx, sm.GetMetric)
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					log.Debug("cannot get cpu: ", err)
					return
				}

				mutex.Lock()
				defer mutex.Unlock()
				tp.point.CPU = data
			}()
		}
	}
}
