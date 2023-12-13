package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func MWCors() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Do something on entering
				// 断言成http的transporter
				ht := tr.(http.Transporter)
				// 在响应中增加特定的头信息
				ht.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
				ht.ReplyHeader().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
				ht.ReplyHeader().Set("Access-Control-Allow-Credentials", "true")
				ht.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,Access-Control-Allow-Credentials,User-Agent,Content-Length,Authorization")
				defer func() {
					// Do something on exiting
				}()
			}
			return handler(ctx, req)
		}
	}
}
