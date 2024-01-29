package data

import (
	"driver/internal/biz"
	"driver/internal/conf"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDriverInterface)

// Data .
type Data struct {
	// TODO wrapped database client
	//Mysql客户端
	MDB *gorm.DB
	//redis客户端
	RDB *redis.Client
	//中间件服务器配置
	cs *conf.Service
}

// NewData .
func NewData(c *conf.Data, cs *conf.Service, logger log.Logger) (*Data, func(), error) {
	//初始化Data数据
	data := &Data{
		cs: cs,
	}

	//初始化RDB
	//1得到一个redis的客户端，连接redis，使用服务的配置项，c就是解析之后的配置信息
	//redis.ParseURL("redis://user:password@localhost:6379/1?dial_timeout=1")
	redisUrl := fmt.Sprintf("redis://%s/1?dial_timeout=%d", c.Redis.Addr, 1)
	options, err := redis.ParseURL(redisUrl)
	if err != nil {

		data.RDB = nil
		log.Fatal(err)
	}
	//new client不会立即连接 需要执行命令才会连接
	data.RDB = redis.NewClient(options) //完成客户端创建
	////ping 测试连接
	//status := rdb.Ping(context.Background())
	//if _, err := status.Result();err != nil {}

	cleanup := func() {
		//匿名接收，清理了Redis连接
		_ = data.RDB.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}

	//二，初始化MDB
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
	return data, cleanup, nil
}
func migrateTable(db *gorm.DB) {
	//自动迁移相关表
	if err := db.AutoMigrate(&biz.Driver{}); err != nil {
		log.Info("Driver table migrate error")
	}
}
