package grpc

import (
	"fmt"

	"github.com/Av4JhG/system-monitoring/config"
	"github.com/Av4JhG/system-monitoring/internal/sm"
	protobuf "github.com/Av4JhG/system-monitoring/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	protobuf.UnimplementedSmServer

	clients sm.NewClienter
	config  config.Config
	log     sm.Logger
}

func newService(log sm.Logger, config config.Config, clients sm.NewClienter) *service {
	return &service{
		clients: clients,
		config:  config,
		log:     log,
	}
}

// GetStats реализует обработку клиентского запроса на получение статистики.
func (s *service) GetStats(req *protobuf.SystemStatisticsRequest, srv protobuf.Sm_GetStatsServer) error {
	s.log.Debug("new client. Every ", req.N, " for ", req.M)

	n := int(req.N)
	m := int(req.M)

	MaxSeconds := s.config.App.MaxSeconds
	if n <= 0 {
		return status.Error(codes.InvalidArgument, "N must be greater than 0 seconds")
	}
	if n > MaxSeconds {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("N must be less than %v seconds", MaxSeconds))
	}
	if m <= 0 {
		return status.Error(codes.InvalidArgument, "M must be greater than 0 seconds")
	}
	if m > MaxSeconds {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("M must be less than %v seconds", MaxSeconds))
	}

	ch, del, err := s.clients.NewClient(sm.ClientData{N: n, M: m})
	if err != nil {
		return status.Error(codes.Unavailable, "service is closing")
	}
	defer del()

L:
	for {
		select {
		case <-srv.Context().Done():
			s.log.Debug("client disconnected")
			break L
		case data, ok := <-ch:
			if !ok {
				break L
			}

			if err := srv.Send(dataToGRPC(data)); err != nil {
				s.log.Debug("unable to send message: ", err)
				break L
			}
		}
	}

	return nil
}

func dataToGRPC(data *sm.Stats) *protobuf.SystemStatistics {
	result := &protobuf.SystemStatistics{}
	result.Time = timestamppb.New(data.Time)

	if data.LoadAvg != nil {
		result.LoadAvg = &protobuf.LoadAvg{
			Load1:  data.LoadAvg.Load1,
			Load5:  data.LoadAvg.Load5,
			Load15: data.LoadAvg.Load15,
		}
	}
	if data.CPU != nil {
		result.Cpu = &protobuf.CPU{
			User:   data.CPU.User,
			System: data.CPU.System,
			Idle:   data.CPU.Idle,
		}
	}
	if data.LoadDisks != nil {
		result.LoadDisks = make([]*protobuf.LoadDisk, 0, len(data.LoadDisks))
		for _, diskData := range data.LoadDisks {
			result.LoadDisks = append(result.LoadDisks, &protobuf.LoadDisk{
				Name:    diskData.Name,
				Tps:     diskData.Tps,
				KBRead:  diskData.KBRead,
				KBWrite: diskData.KBWrite,
			})
		}
	}
	if data.UsedFS != nil {
		result.UsedFs = make([]*protobuf.UsedFS, 0, len(data.UsedFS))
		for _, fsData := range data.UsedFS {
			result.UsedFs = append(result.UsedFs, &protobuf.UsedFS{
				Path:      fsData.Path,
				UsedSpace: fsData.UsedSpace,
				UsedInode: fsData.UsedInode,
			})
		}
	}
	return result
}
