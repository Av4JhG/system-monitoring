package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/sm"
	protobuf "github.com/Av4JhG/system-monitoring/pb"
	"google.golang.org/grpc"
)

type grpcServer struct {
	mutex  *sync.Mutex
	srv    *grpc.Server
	config config.Config
	log    sm.Logger
}

// NewServer возвращает gRPC сервер.
func NewServer(log sm.Logger, config config.Config) sm.GRPCServer {
	return &grpcServer{
		mutex:  &sync.Mutex{},
		config: config,
		log:    log,
	}
}

func (g *grpcServer) Start(addr string, clients sm.NewClienter) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.mutex.Lock()
	g.srv = grpc.NewServer()
	g.mutex.Unlock()

	protobuf.RegisterSmServer(g.srv, newService(g.log, g.config, clients))

	g.log.Debug("starting grpc server on ", addr)

	return g.srv.Serve(lis)
}

func (g *grpcServer) Stop(ctx context.Context) {
	stopped := make(chan interface{})
	go func() {
		g.mutex.Lock()
		defer g.mutex.Unlock()

		g.srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		g.srv.Stop()
	case <-stopped:
		g.srv.Stop()
	}

	g.log.Debug("grpc server is stopped")
}
