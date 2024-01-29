package service

import (
	"context"
	"driver/internal/biz"
	"time"

	pb "driver/api/driver"
)

type DriverService struct {
	pb.UnimplementedDriverServer
	bz *biz.DriverBiz
}

func NewDriverService(bz *biz.DriverBiz) *DriverService {
	return &DriverService{bz: bz}
}

func (s *DriverService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCodeResp, error) {
	code, err := s.bz.GetVerifyCode(ctx, req.Telephone)
	if err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &pb.GetVerifyCodeResp{
		Code:           0,
		Message:        "Success",
		VerifyCode:     code,
		VerifyCodeTime: time.Now().Unix(),
		VerifyCodeLife: 2 * 30 * 24 * 3600,
	}, nil
}
