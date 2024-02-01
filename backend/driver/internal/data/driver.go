package data

import (
	"context"
	"database/sql"
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
func (dd *DriverData) InitDriverInfo(ctx context.Context, tel string) (*biz.Driver, error) {
	//入库，设置状态为stop
	driverInfo := biz.Driver{}
	driverInfo.Telephone = tel
	driverInfo.Status = sql.NullString{
		String: "stop",
		Valid:  true,
	}
	if err := dd.data.MDB.Create(&driverInfo).Error; err != nil {
		return nil, err
	}

	return &driverInfo, nil

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

// 获取已经存储在redis中的验证码
func (dd *DriverData) GetSavedVerifyCode(ctx context.Context, tel string) (string, error) {
	//strCmd := dd.data.RDB.Get(ctx, "DVC"+tel)
	//code, err := strCmd.Result()
	return dd.data.RDB.Get(ctx, "DVC"+tel).Result()

}

// 存储biz层生成的jwtToken 到数据表
func (dd *DriverData) SaveToken(ctx context.Context, tel, token string) error {
	driver := &biz.Driver{}
	//先获取司机信息
	if err := dd.data.MDB.Where("telephone=?", tel).First(&driver).Error; err != nil {
		return err
	}
	//再更新司机信息
	driver.Token = sql.NullString{
		String: token,
		Valid:  true,
	}
	if err := dd.data.MDB.Save(&driver).Error; err != nil {
		return err
	}
	return nil
}
