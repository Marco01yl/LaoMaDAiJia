package biz

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"strconv"
	"valuation/api/mapService"
)

type PriceRule struct {
	gorm.Model
	PriceRuleWork
}

type PriceRuleWork struct {
	CityID      uint  `gorm:"" json:"city_id"`
	StartFee    int64 `gorm:"" json:"start_fee"`
	DistanceFee int64 `gorm:"" json:"distance_fee"`
	DurationFee int64 `gorm:"" json:"duration_fee"`
	StartAt     int   `gorm:"type:int" json:"start_at"` // 0 [0
	EndAt       int   `gorm:"type:int" json:"end_at"`   // 7 0)
}

// 定义操作PriceRule的接口
type PriceRuleInterface interface {
	GetRule(cityId uint, curr int) (*PriceRule, error)
}

type ValuationBiz struct {
	Pri PriceRuleInterface
}

func NewValuationBiz(pri PriceRuleInterface) *ValuationBiz {
	return &ValuationBiz{pri}
}

func (vb *ValuationBiz) GetDrivingInfo(ctx context.Context, origin, destination string) (distance string, duration string, err error) {
	//一、发出请求  grpc
	//一、获取consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return
	}
	//二、获取consul 发现 工具
	dis := consul.New(consulClient) //discover

	endPoint := "discovery:///map"
	//1连接gRPC服务
	conn, err := grpc.DialInsecure(
		context.Background(),
		//grpc.WithEndpoint("localhost:9000"), //grpc服务地址
		grpc.WithEndpoint(endPoint), //目标服务的名字
		//使用服务发现
		grpc.WithDiscovery(dis),
		// 中间件
		grpc.WithMiddleware(
			// tracing 的客户端中间件
			tracing.Client(),
		),
	)
	if err != nil {
		return
	}
	//*******关闭conn。*********
	defer func() {
		_ = conn.Close()
	}()
	//二、发送获取距离和时长//需要在customer服务中，使用map服务中的.proto文件，
	//生成客户端（stub存根)代码，才可以完成grpc的远程调用
	////因此需要拷贝map.proto并使用kratos proto client来生成
	client := mapService.NewMapServiceClient(conn)
	reply, err := client.GetDrivinginfo(context.Background(), &mapService.GetDrivingInfoReq{
		Origin:      origin,
		Destination: destination,
	})
	if err != nil {
		return "", "", nil
	}

	//三、返回正确信息
	distance, duration = reply.Distance, reply.Duration
	return

}

func (vb *ValuationBiz) GetPrice(ctx context.Context, distance, duration string, cityId uint, curr int) (int64, error) {

	//一、获取规则
	rule, err := vb.Pri.GetRule(cityId, curr)
	if err != nil {
		return 0, err
	}
	// 二，将距离和时长，转换为int64
	distanceInt64, err := strconv.ParseInt(distance, 10, 64)
	if err != nil {
		return 0, err
	}
	durationInt64, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return 0, err
	}
	//三，基于rule计算
	distanceInt64 /= 1000
	durationInt64 /= 60
	var startDistance int64 = 5
	totalPrice := rule.StartFee +
		rule.DistanceFee*(distanceInt64-startDistance) +
		rule.DurationFee*(durationInt64)
	return totalPrice, nil

}
