// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/server.proto

package proto

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
	Server_GetAllMetrics_FullMethodName       = "/collectingmetrics.Server/GetAllMetrics"
	Server_SetMetrics_FullMethodName          = "/collectingmetrics.Server/SetMetrics"
	Server_SetMetricsEncrypted_FullMethodName = "/collectingmetrics.Server/SetMetricsEncrypted"
)

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerClient interface {
	GetAllMetrics(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*String, error)
	SetMetrics(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*EmptyMessage, error)
	SetMetricsEncrypted(ctx context.Context, in *EncryptedMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type serverClient struct {
	cc grpc.ClientConnInterface
}

func NewServerClient(cc grpc.ClientConnInterface) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) GetAllMetrics(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, Server_GetAllMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) SetMetrics(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Server_SetMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) SetMetricsEncrypted(ctx context.Context, in *EncryptedMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Server_SetMetricsEncrypted_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
// All implementations must embed UnimplementedServerServer
// for forward compatibility
type ServerServer interface {
	GetAllMetrics(context.Context, *EmptyMessage) (*String, error)
	SetMetrics(context.Context, *Metrics) (*EmptyMessage, error)
	SetMetricsEncrypted(context.Context, *EncryptedMessage) (*EmptyMessage, error)
	mustEmbedUnimplementedServerServer()
}

// UnimplementedServerServer must be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (UnimplementedServerServer) GetAllMetrics(context.Context, *EmptyMessage) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMetrics not implemented")
}
func (UnimplementedServerServer) SetMetrics(context.Context, *Metrics) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMetrics not implemented")
}
func (UnimplementedServerServer) SetMetricsEncrypted(context.Context, *EncryptedMessage) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMetricsEncrypted not implemented")
}
func (UnimplementedServerServer) mustEmbedUnimplementedServerServer() {}

// UnsafeServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerServer will
// result in compilation errors.
type UnsafeServerServer interface {
	mustEmbedUnimplementedServerServer()
}

func RegisterServerServer(s grpc.ServiceRegistrar, srv ServerServer) {
	s.RegisterService(&Server_ServiceDesc, srv)
}

func _Server_GetAllMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetAllMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetAllMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetAllMetrics(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_SetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Metrics)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).SetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_SetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).SetMetrics(ctx, req.(*Metrics))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_SetMetricsEncrypted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncryptedMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).SetMetricsEncrypted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_SetMetricsEncrypted_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).SetMetricsEncrypted(ctx, req.(*EncryptedMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Server_ServiceDesc is the grpc.ServiceDesc for Server service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Server_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "collectingmetrics.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllMetrics",
			Handler:    _Server_GetAllMetrics_Handler,
		},
		{
			MethodName: "SetMetrics",
			Handler:    _Server_SetMetrics_Handler,
		},
		{
			MethodName: "SetMetricsEncrypted",
			Handler:    _Server_SetMetricsEncrypted_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/server.proto",
}