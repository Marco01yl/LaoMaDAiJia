package service

import (
	"context"

	pb "map/api/mapService"
)

type MapServiceService struct {
	pb.UnimplementedMapServiceServer
}

func NewMapServiceService() *MapServiceService {
	return &MapServiceService{}
}

func (s *MapServiceService) GetDrivinginfo(ctx context.Context, req *pb.GetDrivingInfoReq) (*pb.GetDrivingInfoReply, error) {
	return &pb.GetDrivingInfoReply{}, nil
}
