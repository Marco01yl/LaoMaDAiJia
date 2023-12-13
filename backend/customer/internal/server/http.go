package server

import (
	"context"
	"customer/api/customer"
	v1 "customer/api/helloworld/v1"
	"customer/internal/biz"
	"customer/internal/conf"
	"customer/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwt2 "github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	CustomerService *service.CustomerService,
	greeter *service.GreeterService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			//添加自己设置的中间件

			// CORS，全部的请求（响应）都使用该中间件
			selector.Server(MWCors()).Match(func(ctx context.Context, operation string) bool {
				return true
			}).Build(),
			//jwt相关中间件
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(biz.CustomerSecret), nil
				}), customerJWT(CustomerService)).Match(func(ctx context.Context, operation string) bool {
				//根据自己的需求完成是否启用该中间件的校验工作
				noJWT := map[string]struct{}{ //struct结构提相当于key不重复的集合类型
					"/api.customer.Customer/Login":         {},
					"/api.customer.Customer/GetVerifyCode": {},
					//"/api.customer.Customer/Logout":        {},
				}
				if _, exists := noJWT[operation]; exists {
					return false
				}
				return true
			}).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	//注册customer的http服务
	customer.RegisterCustomerHTTPServer(srv, CustomerService)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
