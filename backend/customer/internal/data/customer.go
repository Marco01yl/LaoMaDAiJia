package data

import (
	"context"
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
