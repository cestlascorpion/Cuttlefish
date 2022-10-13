// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/cuttlefish.proto

package cuttlefish

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CuttlefishClient is the client API for Cuttlefish service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CuttlefishClient interface {
}

type cuttlefishClient struct {
	cc grpc.ClientConnInterface
}

func NewCuttlefishClient(cc grpc.ClientConnInterface) CuttlefishClient {
	return &cuttlefishClient{cc}
}

// CuttlefishServer is the server API for Cuttlefish service.
// All implementations must embed UnimplementedCuttlefishServer
// for forward compatibility
type CuttlefishServer interface {
	mustEmbedUnimplementedCuttlefishServer()
}

// UnimplementedCuttlefishServer must be embedded to have forward compatible implementations.
type UnimplementedCuttlefishServer struct {
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

// Cuttlefish_ServiceDesc is the grpc.ServiceDesc for Cuttlefish service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cuttlefish_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Cuttlefish.cuttlefish",
	HandlerType: (*CuttlefishServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "proto/cuttlefish.proto",
}
