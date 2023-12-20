package biz

import (
	"context"
	"customer/api/valuation"
	"database/sql"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

// 业务逻辑部分
type CustomerWork struct {
	Telephone string         `gorm:"type:varchar(15);uniqueIndex;" json:"telephone"`
	Name      sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"name"`
	Email     sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"email"`
	Wechat    sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"wechat"`
	City      uint           `gorm:"index;" json:"city_id"`
}

// Token部分
type CustomerToken struct {
	Token         string       `gorm:"type:varchar(4095);" json:"token"`
	TokenCreateAt sql.NullTime `gorm:"" json:"token_created_at"`
}

// 基础字段部分
type Customer struct {
	//利用gorm嵌入基本的四个字段
	gorm.Model
	//我们需要的业务逻辑
	CustomerWork
	//Token部分
	CustomerToken
}
type CustomerBiz struct {
}

func NewCustomerBiz() *CustomerBiz {
	return &CustomerBiz{}
}

func (cb *CustomerBiz) GetEstimatePrice(origin, destination string) (int64, error) {
	// 一，grpc 获取
	// 1.获取consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	// 2.获取服务发现管理器
	dis := consul.New(consulClient)
	if err != nil {
		return 0, err
	}
	// 2.1,连接目标grpc服务器
	endpoint := "discovery:///Valuation"
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint), // 目标服务的名字
		// 使用服务发现
		grpc.WithDiscovery(dis),
	)
	if err != nil {
		return 0, nil
	}
	//关闭
	defer func() {
		_ = conn.Close()
	}()

	// 2.2,发送获取验证码请求
	client := valuation.NewValuationClient(conn)
	reply, err := client.GetEstimatePrice(context.Background(), &valuation.GetEstimatePriceReq{
		Origin:      origin,
		Destination: destination,
	})
	if err != nil {
		return 0, err
	}

	return reply.Price, nil
}

const CustomerSecret = "yoursecretkey" //加密用字符串要严格保存在服务器端
const CustomerDuration = 2 * 30 * 24 * 3600
