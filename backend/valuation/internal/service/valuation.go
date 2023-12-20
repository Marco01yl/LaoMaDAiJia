package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"valuation/internal/biz"

	pb "valuation/api/valuation"
)

type ValuationService struct {
	pb.UnimplementedValuationServer
	vbBiz *biz.ValuationBiz //从biz包中引用业务实现
}

func NewValuationService(vb *biz.ValuationBiz) *ValuationService {
	return &ValuationService{vbBiz: vb}
}

func (s *ValuationService) GetEstimatePrice(ctx context.Context, req *pb.GetEstimatePriceReq) (*pb.GetEstimatePriceReply, error) {

	//利用biz内的方法计算除距离和时长
	distance, duration, err := s.vbBiz.GetDrivingInfo(ctx, req.Origin, req.Destination)
	if err != nil {
		return nil, errors.New(200, "MAP ERROR", "get driving info error")
	}
	//得到费用
	price, err := s.vbBiz.GetPrice(ctx, distance, duration, 1, 10)
	if err != nil {
		return nil, errors.New(200, "PRICE ERROR", "get price error")
	}

	return &pb.GetEstimatePriceReply{
		Origin:      req.Origin,
		Destination: req.Destination,
		Price:       price,
	}, nil
}
