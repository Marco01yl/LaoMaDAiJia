package biz

import (
	"database/sql"
	"gorm.io/gorm"
)

// 业务逻辑部分
type CustomerWork struct {
	Telephone string         `gorm:"type:varchar(15);uniqueIndex;" json:"telephone"`
	Name      sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"name"`
	Email     sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"email"`
	Wechat    sql.NullString `gorm:"type:varchar(255);uniqueIndex;" json:"wechat"`
	City      uint           `gorm:"index;" json:"city_id"`
}

// Token部分
type CustomerToken struct {
	Token         string       `gorm:"type:varchar(4095);" json:"token"`
	TokenCreateAt sql.NullTime `gorm:"" json:"token_created_at"`
}

// 基础字段部分
type Customer struct {
	//利用gorm嵌入基本的四个字段
	gorm.Model
	//我们需要的业务逻辑
	CustomerWork
	//Token部分
	CustomerToken
}

const CustomerSecret = "yoursecretkey" //加密用字符串要严格保存在服务器端
const CustomerDuration = 2 * 30 * 24 * 3600
