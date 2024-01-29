package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"regexp"
)

type DriverBiz struct {
	di DriverInterface
}

type DriverInterface interface {
	GetVerifyCode(context.Context, string) (string, error)
}

func NewDriverBiz(di DriverInterface) *DriverBiz {
	return &DriverBiz{di: di}
}

func (db *DriverBiz) GetVerifyCode(ctx context.Context, tel string) (string, error) {
	//一、校验手机号
	pattern := `^(13\d|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18\d|19[0-35-9])\d{8}$`
	regexpPattern := regexp.MustCompile(pattern)
	if !regexpPattern.MatchString(tel) {
		return "", errors.New(200, "Driver", "Driver telephone error")
	}
	return db.di.GetVerifyCode(ctx, tel)
}
