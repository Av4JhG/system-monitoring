package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	conf "github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/clients"
	"github.com/Av4JhG/system-monitoring/internal/collector"
	"github.com/Av4JhG/system-monitoring/internal/cpu"
	"github.com/Av4JhG/system-monitoring/internal/grpc"
	loadavg "github.com/Av4JhG/system-monitoring/internal/load_avg"
	loaddisks "github.com/Av4JhG/system-monitoring/internal/load_disks"
	"github.com/Av4JhG/system-monitoring/internal/logger"
	"github.com/Av4JhG/system-monitoring/internal/sm"
	usedfs "github.com/Av4JhG/system-monitoring/internal/used_fs"
	"github.com/benbjohnson/clock"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	if isVersionCommand() {
		printVersion()
		os.Exit(0)
	}

	mainCtx, cancel := context.WithCancel(context.Background())

	go watchSignals(mainCtx, cancel)

	config, err := conf.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("starting system monitor")
	logg.Info("time to keep metrics: ", config.App.MaxSeconds, " seconds")

	stopper := newServiceStopper()

	collectors := sm.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toClientsCh := make(sm.CollectorToClientsCh, 1)

	clientsService := clients.NewClients(logg, clock.New())
	clientsService.Start(mainCtx, toClientsCh)
	stopper.add(clientsService.Stop)

	collectorService := collector.NewCollector(logg, config)
	collectorService.Start(mainCtx, collectors, toClientsCh)
	stopper.add(collectorService.Stop)

	grpcServer := grpc.NewServer(logg, config)
	go func() {
		err := grpcServer.Start(":"+config.Server.Port, clientsService)
		if err != nil {
			logg.Error("server start error: ", err)
			cancel()
			return
		}
	}()
	stopper.add(grpcServer.Stop)

	logg.Info("system monitor is running...")

	<-mainCtx.Done()

	logg.Info("stopping system monitor...")
	stopper.stop()
	logg.Info("system monitor is stopped")
}

func watchSignals(mainCtx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-mainCtx.Done():
	case <-signals:
	}
	cancel()
}
