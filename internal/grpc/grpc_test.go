package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	conf "github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/mocks"
)

func TestGRPCStartStop(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, _ := conf.NewConfig("")

	log := new(mocks.MockLogger)
	log.On("Debug", "grpc server is stopped")
	log.On("Debug", "starting grpc server on ", mock.Anything)

	clientsService := new(mocks.MockNewClienter)

	grpcServer := NewServer(log, config)
	go func() {
		err := grpcServer.Start(":"+config.Server.Port, clientsService)
		require.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond)

	stopCtx := context.Background()
	grpcServer.Stop(stopCtx)

	log.AssertExpectations(t)
}
