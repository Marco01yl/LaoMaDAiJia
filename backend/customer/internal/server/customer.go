package server

import (
	"context"
	"customer/internal/service"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"strings"
)

func customerJWT(customerService *service.CustomerService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 一、获取这个jwt中的id
			claims, ok := jwt.FromContext(ctx) // 得到的claims是interface类型，之后需断言
			if !ok {
				//没有获取到claims
				return nil, errors.Unauthorized("UNAUTHORIZED", "claims not found")
			}
			//1.2断言claims
			claimsMap := claims.(jwt2.MapClaims)
			//map中jti字段就是浏览器存储格式json中的jwt的id
			id := claimsMap["jti"]
			//二、获取对应id的token
			token, err := customerService.CD.GetToken(id)
			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED", "customer not found")
			}
			//比对数据表中的token与请求的token是否一致
			// 获取请求头
			header, _ := transport.FromServerContext(ctx)
			// 从header获取token
			auths := strings.SplitN(header.RequestHeader().Get("Authorization"), " ", 2)
			jwtToken := auths[1]
			// 比较请求中的token与数据表中获取的token是否一致
			if jwtToken != token {
				return nil, errors.Unauthorized("UNAUTHORIZED", "token was updated")
			}
			//四、校验通过，放行，交由下个中间件（handler） 处理
			return handler(ctx, req)
		}
	}
}
