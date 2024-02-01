package service

import (
	"context"
	"driver/internal/biz"
	"log"
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
func (s *DriverService) SubmitPhone(ctx context.Context, req *pb.SubmitPhoneReq) (*pb.SubmitPhoneResp, error) {
	//首先校验验证码（略）
	//司机是否已经注册的校验（略）
	//司机是否在黑名单（略）

	//将司机信息入库,并设置状态为stop，暂停使用（核心逻辑）需要在biz中增加功能，上面的功能也在biz中实现
	driver, err := s.bz.InitDriverInfo(ctx, req.Telephone)
	if err != nil {
		return &pb.SubmitPhoneResp{
			Code:    1,
			Message: "司机号码提交失败",
		}, nil
	}
	return &pb.SubmitPhoneResp{
		Code:    0,
		Message: "司机号码提交成功",
		Status:  driver.Status.String,
	}, nil
}
func (s *DriverService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {

	//由biz层完成业务逻辑处理
	token, err := s.bz.CheckLogin(ctx, req.Telephone, req.VerifyCode)
	if err != nil {
		//具体原因可通过log记录
		log.Println(err)
		return &pb.LoginResp{
			Code:    1,
			Message: "司机登录失败t",
		}, nil
	}
	return &pb.LoginResp{
		Code:          0,
		Message:       "司机登录成功",
		Token:         token,
		TokenCreateAt: time.Now().Unix(),
		TokenLife:     biz.TokenLifetime,
	}, nil
}
