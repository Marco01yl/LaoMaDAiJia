package data

import (
	"customer/internal/conf"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
// NewCustomerData手动添加
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCustomerData)

// Data
type Data struct {
	// TODO wrapped database client
	//手动加入初始化Redis的客户端
	RDB *redis.Client
	MDB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	//初始化Data数据
	data := &Data{}
	//初始化RDB
	//1得到一个redis的客户端，连接redis，使用服务的配置项，c就是解析之后的配置信息
	//redis.ParseURL("redis://user:password@localhost:6379/1?dial_timeout=1")
	redisUrl := fmt.Sprintf("redis://%s/1?dial_timeout=%d", c.Redis.Addr, 1)
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		data.RDB = nil
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

	return data, cleanup, nil
}
