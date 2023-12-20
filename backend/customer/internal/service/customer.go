package service

import (
	"context"
	pb "customer/api/customer"
	verifyCode "customer/api/verifyCode"
	"customer/internal/biz"
	"customer/internal/data"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/consul/api"
	"regexp"
	"time"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	CD   *data.CustomerData
	Cbiz *biz.CustomerBiz
}

func NewCustomerService(cd *data.CustomerData, cb *biz.CustomerBiz) *CustomerService {
	return &CustomerService{
		CD:   cd,
		Cbiz: cb,
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

	//一、获取consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}
	//二、获取consul 发现 工具
	dis := consul.New(consulClient) //discover

	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())
	//selector.SetGlobalSelector(p2c.NewBuilder())
	endPoint := "discovery:///verifyCode"
	//1连接gRPC服务
	conn, err := grpc.DialInsecure(
		context.Background(),
		//grpc.WithEndpoint("localhost:9000"), //grpc服务地址
		grpc.WithEndpoint(endPoint), //目标服务的名字
		//使用服务发现
		grpc.WithDiscovery(dis),
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
	if err := s.CD.SetVerifyCode(req.Telephone, reply.Code, life); err != nil {
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
	code := s.CD.GetVerifyCode(req.Telephone)
	if code == "" || code != req.VerifyCode {
		return &pb.LoginResp{
			Code:    1,
			Message: "校验不正确",
		}, nil
	}
	//二、判断号码是否已经注册
	customer, err := s.CD.GetCustomerByTelephone(req.Telephone)
	if err != nil {
		return &pb.LoginResp{
			Code:    1,
			Message: "顾客信息获取错误",
		}, nil
	}
	//三、设置Token， jwt-token
	//const secret = "yoursecretkey" //加密用字符串要严格保存在服务器端
	//const duration = 2 * 30 * 24 * 3600

	token, err := s.CD.GenerateTokenAndSave(customer, biz.CustomerDuration*time.Second, biz.CustomerSecret)
	if err != nil {
		return &pb.LoginResp{
			Code:    1,
			Message: "Token created failed.",
			Token:   token,
		}, nil
	}
	//四、响应token
	return &pb.LoginResp{
		Code:          0,
		Message:       "Token created successed.",
		Token:         token,
		TokenCreateAt: time.Now().Unix(),
		//TokenLife:     2 * 30 * 24 * 3600, 可已设置为常量
		TokenLife: biz.CustomerDuration,
	}, nil
}

// Code:          0,
// Message:       "Token created successed.",
// Token:         token,
// TokenCreateAt: time.Now().Unix(),
// //TokenLife:     2 * 30 * 24 * 3600, 可已设置为常量
// TokenLife: duration,
func (s *CustomerService) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutResp, error) {

	//获得用户的信息
	// 一、获取这个jwt中的id
	claims, _ := jwt.FromContext(ctx) // 得到的claims是interface类型，之后需断言
	//1.2断言claims
	claimsMap := claims.(jwt2.MapClaims)
	//map中jti字段就是浏览器存储格式json中的jwt的id
	//id := claimsMap["jti"]
	//删除用户的token
	if err := s.CD.DleToken(claimsMap["jti"]); err != nil {
		return &pb.LogoutResp{
			Code:    1,
			Message: "Token 删除失败",
		}, nil
	}
	//成功，响应
	return &pb.LogoutResp{
		Code:    0,
		Message: "logout success",
	}, nil
}
func (s *CustomerService) EstimatePrice(ctx context.Context, req *pb.EstimatePriceReq) (*pb.EstimatePriceResp, error) {
	price, err := s.Cbiz.GetEstimatePrice(req.Origin, req.Destination)
	if err != nil {
		return &pb.EstimatePriceResp{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &pb.EstimatePriceResp{
		Code:        0,
		Message:     "got EstimatePrice",
		Origin:      req.Origin,
		Destination: req.Destination,
		Price:       price,
	}, nil
}
