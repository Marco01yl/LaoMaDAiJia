package biz

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"regexp"
)

type DriverBiz struct {
	di DriverInterface
}

type DriverInterface interface {
	GetVerifyCode(context.Context, string) (string, error)
	InitDriverInfo(context.Context, string) (*Driver, error)
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

// 司机表模型
type Driver struct {
	// 基础模型
	gorm.Model
	// 业务模型
	DriverWork
	// 关联部分
}

// 司机的业务模型
type DriverWork struct {
	Telephone     string         `gorm:"type:varchar(16);uniqueIndex;" json:"telephone"`
	Token         sql.NullString `gorm:"type:varchar(2047);" json:"token"`
	Name          sql.NullString `gorm:"type:varchar(255);index;" json:"name"`
	Status        sql.NullString `gorm:"type:enum('out', 'in', 'listen', 'stop');" json:"status"`
	IdNumber      sql.NullString `gorm:"type:char(18);uniqueIndex;" json:"id_number"`
	IdImageA      sql.NullString `gorm:"type:varchar(255);" json:"id_image_a"`
	LicenseImageA sql.NullString `gorm:"type:varchar(255);" json:"license_image_a"`
	LicenseImageB sql.NullString `gorm:"type:varchar(255);" json:"license_image_b"`
	DistinctCode  sql.NullString `gorm:"type:varchar(16);index;" json:"distinct_code"`
	TelephoneBak  sql.NullString `gorm:"type:varchar(16);index;" json:"telephone_bak"`
	AuditAt       sql.NullTime   `gorm:"index;" json:"audit_at"`
}

// 司机状态常量
const DriverStatusOut = "out"
const DriverStatusIn = "in"
const DriverStatusListen = "listen"
const DriverStatusStop = "stop"

// 将司机信息入库的功能
func (db *DriverBiz) InitDriverInfo(ctx context.Context, tel string) (*Driver, error) {
	//在业务层做一些判断，在数据层进行数据的创建（向上回到service层）
	if tel == "" {
		return nil, errors.New(1, "telephone is empty", "")
	}
	return db.di.InitDriverInfo(ctx, tel)
}
