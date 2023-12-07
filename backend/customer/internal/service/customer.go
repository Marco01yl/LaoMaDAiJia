package service

import (
	"context"
	pb "customer/api/customer"
	verifyCode "customer/api/verifyCode"
	"customer/internal/data"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"regexp"
	"time"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	cd *data.CustomerData
}

func NewCustomerService(cd *data.CustomerData) *CustomerService {
	return &CustomerService{
		cd: cd,
	}
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
	client := verifyCode.NewVerifyCodeClient(conn)
	reply, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   1,
	})
	if err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "get verify-code failed",
		}, nil
	}

	const life = 60 //定义临时缓存时间
	//3、redis的临时存储
	//使用go-redis包完成redis操作
	if err := s.cd.SetVerifyCode(req.Telephone, reply.Code, life); err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "get verify-code failed(set redis)",
		}, nil
	}

	//生成相应
	return &pb.GetVerifyCodeResp{
		Code:               0,
		VerifyCode:         reply.Code,
		VerifyCodeTime:     time.Now().Unix(),
		VerifyCodeLifetime: life,
	}, nil
}

// 手动写入登录的业务逻辑
func (s *CustomerService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	//校验手机和验证码
	//用电话号
	//key := "CVC:" + req.Telephone
	code := s.cd.GetVerifyCode(req.Telephone)
	if code == "" || code != req.VerifyCode {
		return &pb.LoginResp{
			Code:    1,
			Message: "验证码不匹配",
		}, nil
	}
	//二、判断号码是否已经注册
	customer, err := s.cd.GetCustomerByTelephone(req.Telephone)
	if err != nil {
		return &pb.LoginResp{
			Code:    1,
			Message: "顾客信息获取错误",
		}, nil
	}
	//设置Token， iwt-token
	token := s.cd.GenerateTokenAndSave(customer)
	//响应token
	return &pb.LoginResp{
		Code:          0,
		Message:       "Token created.",
		Token:         token.Token,
		TokenCreateAt: token.TokenCreateAt,
		TokenLife:     token.TokenLife,
	}, nil
}
