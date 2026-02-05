package sm

import (
	"context"
	"errors"
	"time"
)

// ErrStopped ошибка, возвращаемая grpc запросу, если приложение останавливается.
var ErrStopped = errors.New("service is stopped")

// Collector представляет сервис, запускающий каждую секунду сбор статистики и ее отправку клиентам.
type Collector interface {
	Start(context.Context, MetricCollectors, chan<- MetricsData)
	Stop(context.Context)
}

// NewClienter представляет интерфейс для подключения новых клиентов.
type NewClienter interface {
	// возвращает канал для получения отсылаемых данных и ф-ия отключения клиента
	NewClient(ClientData) (<-chan *Stats, func(), error)
}

// Clients представляет сервис, хранящий всех подключенных клиентов и отсылающий им статистику.
type Clients interface {
	Start(context.Context, <-chan MetricsData)
	Stop(context.Context)
	NewClienter
}

// CollectorToClientsCh - канал для посекундной передачи накопленных данных сервису клиентов.
type CollectorToClientsCh chan MetricsData

// MetricsData содержит данные, отсылаемые сервису клиентов.
// Отсылается текущая секунда и копия всех собранных данных.
// Отсылается копия, чтобы сервис мог ее обрабатывать, не блокируя мьютекс с собираемыми данными.
type MetricsData struct {
	Time   time.Time
	Points Points
}

// Stats содержит данные, отсылаемые каждому клиенту.
type Stats struct {
	Time      time.Time
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

// Points хранит собранные посекундные наборы метрик.
type Points map[time.Time]*Point

// Point содержит набор метрик (снапшот). За секунду или усредненный.
type Point struct {
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

// MetricCommand - команды для взаимодействия сервиса метрик и коллекторами, собирающими метрики.
type MetricCommand int

const (
	// StartMetric - начать собирать метрики.
	StartMetric MetricCommand = iota
	// StopMetric - остановить сбор метрик.
	StopMetric
	// GetMetric - получить метрики.
	GetMetric
)

// MetricCollectors - набор функций, возвращающих свои метрики. Передается сервису Collector при его создании.
type MetricCollectors struct {
	LoadAvg   LoadAvg
	CPU       CPU
	LoadDisks LoadDisks
	UsedFS    UsedFS
}

// GRPCServer представляет gRPC сервер.
type GRPCServer interface {
	Start(addr string, clients NewClienter) error
	Stop(ctx context.Context)
}

// ClientData - информация, передаваемая из grpc запроса сервису клиентов.
type ClientData struct {
	N int // информация отправляется каждые N секунд
	M int // информация усредняется за M секунд
}

// Logger представляет логгер.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}
