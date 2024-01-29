// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.6
// source: api/driver/driver.proto

package driver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Driver_GetVerifyCode_FullMethodName = "/api.driver.Driver/GetVerifyCode"
)

// DriverClient is the client API for Driver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DriverClient interface {
	// 获取验证码服务
	GetVerifyCode(ctx context.Context, in *GetVerifyCodeReq, opts ...grpc.CallOption) (*GetVerifyCodeResp, error)
}

type driverClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverClient(cc grpc.ClientConnInterface) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) GetVerifyCode(ctx context.Context, in *GetVerifyCodeReq, opts ...grpc.CallOption) (*GetVerifyCodeResp, error) {
	out := new(GetVerifyCodeResp)
	err := c.cc.Invoke(ctx, Driver_GetVerifyCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverServer is the server API for Driver service.
// All implementations must embed UnimplementedDriverServer
// for forward compatibility
type DriverServer interface {
	// 获取验证码服务
	GetVerifyCode(context.Context, *GetVerifyCodeReq) (*GetVerifyCodeResp, error)
	mustEmbedUnimplementedDriverServer()
}

// UnimplementedDriverServer must be embedded to have forward compatible implementations.
type UnimplementedDriverServer struct {
}

func (UnimplementedDriverServer) GetVerifyCode(context.Context, *GetVerifyCodeReq) (*GetVerifyCodeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVerifyCode not implemented")
}
func (UnimplementedDriverServer) mustEmbedUnimplementedDriverServer() {}

// UnsafeDriverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DriverServer will
// result in compilation errors.
type UnsafeDriverServer interface {
	mustEmbedUnimplementedDriverServer()
}

func RegisterDriverServer(s grpc.ServiceRegistrar, srv DriverServer) {
	s.RegisterService(&Driver_ServiceDesc, srv)
}

func _Driver_GetVerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVerifyCodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).GetVerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Driver_GetVerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).GetVerifyCode(ctx, req.(*GetVerifyCodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Driver_ServiceDesc is the grpc.ServiceDesc for Driver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Driver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.driver.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVerifyCode",
			Handler:    _Driver_GetVerifyCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/driver/driver.proto",
}
