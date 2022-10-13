// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/cuttlefish.proto

package cuttlefish

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

// CuttlefishClient is the client API for Cuttlefish service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CuttlefishClient interface {
	GetTentacle(ctx context.Context, in *GetTentacleReq, opts ...grpc.CallOption) (*GetTentacleResp, error)
	SetTentacle(ctx context.Context, in *SetTentacleReq, opts ...grpc.CallOption) (*SetTentacleResp, error)
}

type cuttlefishClient struct {
	cc grpc.ClientConnInterface
}

func NewCuttlefishClient(cc grpc.ClientConnInterface) CuttlefishClient {
	return &cuttlefishClient{cc}
}

func (c *cuttlefishClient) GetTentacle(ctx context.Context, in *GetTentacleReq, opts ...grpc.CallOption) (*GetTentacleResp, error) {
	out := new(GetTentacleResp)
	err := c.cc.Invoke(ctx, "/Cuttlefish.cuttlefish/GetTentacle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cuttlefishClient) SetTentacle(ctx context.Context, in *SetTentacleReq, opts ...grpc.CallOption) (*SetTentacleResp, error) {
	out := new(SetTentacleResp)
	err := c.cc.Invoke(ctx, "/Cuttlefish.cuttlefish/SetTentacle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CuttlefishServer is the server API for Cuttlefish service.
// All implementations must embed UnimplementedCuttlefishServer
// for forward compatibility
type CuttlefishServer interface {
	GetTentacle(context.Context, *GetTentacleReq) (*GetTentacleResp, error)
	SetTentacle(context.Context, *SetTentacleReq) (*SetTentacleResp, error)
	mustEmbedUnimplementedCuttlefishServer()
}

// UnimplementedCuttlefishServer must be embedded to have forward compatible implementations.
type UnimplementedCuttlefishServer struct {
}

func (UnimplementedCuttlefishServer) GetTentacle(context.Context, *GetTentacleReq) (*GetTentacleResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTentacle not implemented")
}
func (UnimplementedCuttlefishServer) SetTentacle(context.Context, *SetTentacleReq) (*SetTentacleResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTentacle not implemented")
}
func (UnimplementedCuttlefishServer) mustEmbedUnimplementedCuttlefishServer() {}

// UnsafeCuttlefishServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CuttlefishServer will
// result in compilation errors.
type UnsafeCuttlefishServer interface {
	mustEmbedUnimplementedCuttlefishServer()
}

func RegisterCuttlefishServer(s grpc.ServiceRegistrar, srv CuttlefishServer) {
	s.RegisterService(&Cuttlefish_ServiceDesc, srv)
}

func _Cuttlefish_GetTentacle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTentacleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CuttlefishServer).GetTentacle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cuttlefish.cuttlefish/GetTentacle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CuttlefishServer).GetTentacle(ctx, req.(*GetTentacleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cuttlefish_SetTentacle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetTentacleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CuttlefishServer).SetTentacle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cuttlefish.cuttlefish/SetTentacle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CuttlefishServer).SetTentacle(ctx, req.(*SetTentacleReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Cuttlefish_ServiceDesc is the grpc.ServiceDesc for Cuttlefish service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cuttlefish_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Cuttlefish.cuttlefish",
	HandlerType: (*CuttlefishServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTentacle",
			Handler:    _Cuttlefish_GetTentacle_Handler,
		},
		{
			MethodName: "SetTentacle",
			Handler:    _Cuttlefish_SetTentacle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/cuttlefish.proto",
}
