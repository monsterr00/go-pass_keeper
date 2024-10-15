package httplayer

import (
	"context"

	config "github.com/monsterr00/go-pass-keeper/configs/server"
	"github.com/monsterr00/go-pass-keeper/internal/server/applayer"
	pb "github.com/monsterr00/go-pass-keeper/internal/server/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DataServer struct {
	pb.UnimplementedKeeperDataServer
	appRepo applayer.App
}

// AddData обрабатывает запросы на добавление/обновление данных.
func (s *DataServer) AddData(ctx context.Context, in *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	var response pb.AddDataResponse

	err := s.appRepo.Add(ctx, in.Data.DataType, in.Data.Data)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "Server: add data error")
	}
	return &response, nil
}

// GetData обрабатывает запросы на чтение данных.
func (s *DataServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var response pb.GetDataResponse
	var err error
	var data pb.Data

	data.DataType = config.TextDataType
	data.Data, err = s.appRepo.Get(ctx, in.Data.DataType, in.Data.Data)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "Server: get data error")
	}

	response.Data = &data
	return &response, nil
}
