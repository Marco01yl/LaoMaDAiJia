package data

import (
	"context"
	"driver/api/verifyCode"
	"driver/internal/biz"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"time"
)

type DriverData struct {
	data *Data
}

func NewDriverInterface(data *Data) biz.DriverInterface {
	return &DriverData{data: data}
}

func (dd *DriverData) GetVerifyCode(ctx context.Context, tel string) (string, error) {
	//grpc请求
	consulConfig := api.DefaultConfig()
	consulConfig.Address = dd.data.cs.Consul.Address
	consulClient, err := api.NewClient(consulConfig)
	discov := consul.New(consulClient)
	if err != nil {
		return "", err
	}
	endpoint := "discovery:///VerifyCode"
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(discov),
	)
	if err != nil {
		return "", err
	}
	//关闭
	defer conn.Close()

	//2.2 发送获取验证码请求
	client := verifyCode.NewVerifyCodeClient(conn)
	reply, err := client.GetVerifyCode(ctx, &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   1,
	})
	if err != nil {
		return "", err
	}
	//redis的临时存储
	status := dd.data.RDB.Set(ctx, "DVC"+tel, reply.Code, 60*time.Second)
	if _, err := status.Result(); err != nil {
		return "", err
	}
	return reply.Code, nil
}
