package data

import (
	"context"
	"customer/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"time"
)

// 定义customer中与数据操作相关的操作
type CustomerData struct {
	data *Data
}

//New方法

func NewCustomerData(data *Data) *CustomerData {
	return &CustomerData{data: data}
}

// 将业务逻辑中设置vcode的方法放到这里
func (cd CustomerData) SetVerifyCode(telephone, code string, ex int) error {

	//业务命令即可相当于开始连接//设置key。 customer-verify-code=>cvc + req.telephone
	status := cd.data.RDB.Set(context.Background(), "CVC:"+telephone, code, time.Duration(ex)*time.Second)
	if _, err := status.Result(); err != nil {
		return err
	}

	return nil
}

// 业务逻辑中校验用户和验证码以及生成token的逻辑
// 获取对应的验证码
func (cd CustomerData) GetVerifyCode(telephone string) string {
	status := cd.data.RDB.Get(context.Background(), "CVC:"+telephone)
	return status.String()
}

// 根据电话获取顾客信息
func (cd CustomerData) GetCustomerByTelephone(telephone string) (*biz.Customer, error) {
	custoemr := &biz.Customer{}
	result := cd.data.MDB.Where("telephone=?").First(custoemr)

	if result.Error == nil && custoemr.ID > 0 {
		//query执行成功，返回cusetomer
		return custoemr, nil
	}
	//不成功
	//记录不存在。创建customer并返回
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {

		custoemr.Telephone = telephone
		resultCreate := cd.data.MDB.Create(custoemr)
		//创建插入成功
		if resultCreate.Error != nil {
			return custoemr, nil
		} else {
			return nil, resultCreate.Error
		}

	}
	//不是记录不存在。不做业务逻辑处理
	return nil, result.Error
}
