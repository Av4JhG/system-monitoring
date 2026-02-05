package collector

import (
	"context"
	"sync"
	"time"

	"github.com/Av4JhG/system-monitoring/internal/sm"
)

// общий код для тестирования всех горутин, отвечающих за сбор конкретных метрик.
func testCollector() (context.Context, sync.Locker, <-chan timePoint, *sm.Point) {
	ctx, cancel := context.WithCancel(context.Background())
	mutex := &sync.Mutex{}
	ch := make(chan timePoint, 1)

	point := &sm.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	return ctx, mutex, ch, point
}
