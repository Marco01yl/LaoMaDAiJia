package data

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"valuation/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"valuation/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewPriceRuleInterface)

// Data .
type Data struct {
	// TODO wrapped database client
	MDB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}
	//初始化MDB
	//连接mysql，使用配置
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		data.MDB = nil
		fmt.Println("mysql started failed")
		log.Fatal(err)
	}
	data.MDB = db
	//开发阶段，用func migrateTable()迁移表结构; 发布阶段表结构稳定不需要再次migrate
	migrateTable(db)
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}

func migrateTable(db *gorm.DB) {

	if err := db.AutoMigrate(&biz.PriceRule{}); err != nil {
		log.Info("price_rule table migrate error,err:", err)
	}

	// 插入一些riceRule的测试数据
	rules := []biz.PriceRule{
		{
			Model: gorm.Model{ID: 1},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    300,
				DistanceFee: 35,
				DurationFee: 10, // 5m
				StartAt:     7,
				EndAt:       23,
			},
		},
		{
			Model: gorm.Model{ID: 2},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    350,
				DistanceFee: 35,
				DurationFee: 10, // 5m
				StartAt:     23,
				EndAt:       24,
			},
		},
		{
			Model: gorm.Model{ID: 3},
			PriceRuleWork: biz.PriceRuleWork{
				CityID:      1,
				StartFee:    400,
				DistanceFee: 35,
				DurationFee: 10, // 5m
				StartAt:     0,
				EndAt:       7,
			},
		},
	}
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(rules)
}
