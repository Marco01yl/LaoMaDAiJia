package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"regexp"

	pb "customer/api/customer"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCodeResp, error) {
	//一、校验手机号
	pattern := regexp.MustCompile(`^(13\d|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18\d|19[0-35-9])\d{8}$`)
	if !pattern.MatchString(req.Telephone) {
		// 如果正则匹配失败，则返回错误
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "wrong number",
		}, nil
	}

	//二、通过验证码服务生成验证码（服务间通信，grpc）//1、连接gRPC服务2、发送获取验证码请求
	//1连接gRPC服务
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("localhost:9000"), //grpc服务地址
	)
	if err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "verify-code service unvalid",
		}, nil
	}
	//*******关闭conn。*********
	defer func() {
		_ = conn.Close()
	}()
	//2、发送获取验证码请求//需要在customer服务中，使用verify-code服务中的.proto文件，
	//生成客户端（stub存根)代码，才可以完成grpc的远程调用
	////因此需要拷贝verifyCode.proto并使用kratos client来生成

	return &pb.GetVerifyCodeResp{}, nil
}