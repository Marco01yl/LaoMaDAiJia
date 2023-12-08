package data

import (
	"context"
	"customer/internal/biz"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v4"
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
	code, _ := status.Result() // status.String()
	return code
}

// 根据电话获取顾客信息
func (cd CustomerData) GetCustomerByTelephone(telephone string) (*biz.Customer, error) {
	customer := &biz.Customer{}
	result := cd.data.MDB.Where("telephone=?", telephone).First(customer)

	if result.Error == nil && customer.ID > 0 {
		//query执行成功，返回customer
		return customer, nil
	}
	//不成功
	//记录不存在。创建customer并返回
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {

		customer.Telephone = telephone
		customer.Email = sql.NullString{Valid: false}
		customer.Wechat = sql.NullString{Valid: false}
		customer.Name = sql.NullString{Valid: false}
		resultCreate := cd.data.MDB.Create(customer)
		//创建插入成功
		if resultCreate.Error != nil {
			return customer, nil
		} else {
			return nil, resultCreate.Error
		}

	}
	//不是记录不存在。不做业务逻辑处理
	return nil, result.Error
}

// 生成token和存储
func (cd CustomerData) GenerateTokenAndSave(c *biz.Customer, duration time.Duration, secret string) (string, error) {
	//一、生成token
	//处理token中的数据 //利用标准的jwt的payload格式
	claims := jwt.RegisteredClaims{
		//签发方
		Issuer: "LaoMaDJ",
		//简单说明
		Subject: "customer-authentication",
		//签发给谁使用
		Audience: []string{"Customer", "others"},
		//有效期至
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		//何时启用
		NotBefore: nil,
		//签发时间
		IssuedAt: jwt.NewNumericDate(time.Now()),
		//可用来存储用户的ID
		ID: fmt.Sprintf("%d", c.ID),
	}
	//利用生成的payload生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//签名token
	secretToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	//二、存储
	c.Token = secretToken
	c.TokenCreateAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if result := cd.data.MDB.Save(c); result.Error != nil {
		return "", result.Error
	}
	//三、操作完毕，返回token
	return secretToken, nil
}
