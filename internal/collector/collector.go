package collector

import (
	"context"
	"sync"
	"time"

	conf "github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/sm"
)

const timeToGetMetric = 950 * time.Millisecond

type collector struct {
	ctx         context.Context // управление остановкой сервиса
	ctxCancel   context.CancelFunc
	stoppedCh   chan interface{}
	mutex       *sync.Mutex
	points      sm.Points          // собираемые данные
	workerChans []chan<- timePoint // каналы горутин, ответственных за сбор конкретных метрик
	config      conf.Config
	collectors  sm.MetricCollectors // функции возвращающие конкретные метрики
	toClientsCh chan<- sm.MetricsData
	log         sm.Logger
}

// информация, отправляемая горутинам, ответственным за сбор конкретных метрик.
type timePoint struct {
	time  time.Time // за какую секунду метрика
	point *sm.Point // структура, в которую складываются метрики
}

// NewCollector возвращает сервис сбора метрик.
func NewCollector(log sm.Logger, config conf.Config) sm.Collector {
	return &collector{
		config: config,
		log:    log,
	}
}

func (c *collector) Start(ctx context.Context, collectors sm.MetricCollectors, toClientsCh chan<- sm.MetricsData) {
	c.collectors = collectors
	c.toClientsCh = toClientsCh

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.stoppedCh = make(chan interface{})
	c.mutex = &sync.Mutex{}
	c.points = make(sm.Points)
	c.workerChans = nil

	mountedCh := make(chan interface{})
	go c.mountMetrics(ctx, mountedCh)

	select {
	case <-ctx.Done():
		// стартануть не успели, помечаем сервис как остановленный
		close(c.stoppedCh)
		return
	case <-mountedCh:
		go c.work()
	}
}

func (c *collector) Stop(ctx context.Context) {
	c.ctxCancel()

	select {
	case <-ctx.Done():
		return
	case <-c.stoppedCh:
	}

	unmountedCh := make(chan interface{})
	go c.unmountMetrics(ctx, unmountedCh)

	select {
	case <-ctx.Done():
		return
	case <-unmountedCh:
		c.log.Debug("collector is stopped")
	}
}

func (c *collector) mountMetrics(startCtx context.Context, mountedCh chan interface{}) {
	wg := &sync.WaitGroup{}

	if c.config.Metric.LoadAvg {
		go loadavgCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.LoadAvg, c.log)
	}
	if c.config.Metric.CPU {
		wg.Add(1)
		go c.mountCPU(startCtx, wg)
	}
	if c.config.Metric.LoadDisks {
		wg.Add(1)
		go c.mountLoadDisks(startCtx, wg)
	}
	if c.config.Metric.UsedFS {
		wg.Add(1)
		go c.mountUsedFS(startCtx, wg)
	}

	wg.Wait()
	close(mountedCh)
}

func (c *collector) mountCPU(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.CPU(startCtx, sm.StartMetric)
	if err != nil {
		c.log.Debug("cannot start collect the cpu metric: ", err)
		return
	}
	go cpuCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.CPU, c.log)
}

func (c *collector) mountLoadDisks(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.LoadDisks(startCtx, sm.StartMetric)
	if err != nil {
		c.log.Debug("cannot start collect the load disks metric: ", err)
	}
	go loadDisksCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.LoadDisks, c.log)
}

func (c *collector) mountUsedFS(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.UsedFS(startCtx, sm.StartMetric)
	if err != nil {
		c.log.Debug("cannot start collect the used fs metric: ", err)
	}
	go usedFSCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.UsedFS, c.log)
}

func (c *collector) unmountMetrics(stopCtx context.Context, unmountedCh chan interface{}) {
	wg := &sync.WaitGroup{}

	if c.config.Metric.LoadDisks {
		wg.Add(1)
		go c.unmountLoadDisks(stopCtx, wg)
	}
	if c.config.Metric.UsedFS {
		wg.Add(1)
		go c.unmountUsedFS(stopCtx, wg)
	}

	wg.Wait()
	close(unmountedCh)
}

func (c *collector) unmountLoadDisks(stopCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.LoadDisks(stopCtx, sm.StopMetric)
	if err != nil {
		c.log.Debug("cannot stop load disks: ", err)
	}
}

func (c *collector) unmountUsedFS(stopCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.UsedFS(stopCtx, sm.StopMetric)
	if err != nil {
		c.log.Debug("cannot stop used fs: ", err)
	}
}

func (c *collector) newWorkerChan() chan timePoint {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ch := make(chan timePoint, 1)
	c.workerChans = append(c.workerChans, ch)

	return ch
}

func (c *collector) work() {
	defer close(c.stoppedCh)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		select {
		case <-c.ctx.Done():
			return
		case now := <-ticker.C:
			c.processTick(now.Truncate(time.Second))
		}
	}
}

func (c *collector) processTick(now time.Time) {
	c.log.Debug("tick ", now)

	// добавляется новая точка для статистики за эту секунду
	point := c.addPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range c.workerChans {
		tp := timePoint{
			time:  now,
			point: point,
		}
		select {
		case ch <- tp:
		default:
		}
	}

	// устаревшие точки удаляются
	c.cleanPoints(now)

	data := sm.MetricsData{
		Time:   now,
		Points: c.cloneOldPoints(now),
	}

	select {
	case c.toClientsCh <- data:
	default:
	}
}

func (c *collector) addPoint(now time.Time) *sm.Point {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	point := &sm.Point{}
	c.points[now] = point

	return point
}

func (c *collector) cleanPoints(now time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	limit := now.Add(-time.Duration(c.config.App.MaxSeconds) * time.Second)

	for key := range c.points {
		if key.Before(limit) {
			delete(c.points, key)
		}
	}
}

func (c *collector) cloneOldPoints(now time.Time) sm.Points {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make(sm.Points, len(c.points))
	for key, point := range c.points {
		if key.Equal(now) {
			continue
		}
		newPoint := *point
		result[key] = &newPoint
	}

	return result
}
