// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: v1/looncan.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	LooncanService_List_FullMethodName          = "/microservice.LooncanService/List"
	LooncanService_ListForParent_FullMethodName = "/microservice.LooncanService/ListForParent"
)

// LooncanServiceClient is the client API for LooncanService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LooncanServiceClient interface {
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListLooncanResponse, error)
	ListForParent(ctx context.Context, in *ListLooncanForParentRequest, opts ...grpc.CallOption) (*ListLooncanResponse, error)
}

type looncanServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLooncanServiceClient(cc grpc.ClientConnInterface) LooncanServiceClient {
	return &looncanServiceClient{cc}
}

func (c *looncanServiceClient) List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListLooncanResponse, error) {
	out := new(ListLooncanResponse)
	err := c.cc.Invoke(ctx, LooncanService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *looncanServiceClient) ListForParent(ctx context.Context, in *ListLooncanForParentRequest, opts ...grpc.CallOption) (*ListLooncanResponse, error) {
	out := new(ListLooncanResponse)
	err := c.cc.Invoke(ctx, LooncanService_ListForParent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LooncanServiceServer is the server API for LooncanService service.
// All implementations must embed UnimplementedLooncanServiceServer
// for forward compatibility
type LooncanServiceServer interface {
	List(context.Context, *emptypb.Empty) (*ListLooncanResponse, error)
	ListForParent(context.Context, *ListLooncanForParentRequest) (*ListLooncanResponse, error)
	mustEmbedUnimplementedLooncanServiceServer()
}

// UnimplementedLooncanServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLooncanServiceServer struct {
}

func (UnimplementedLooncanServiceServer) List(context.Context, *emptypb.Empty) (*ListLooncanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedLooncanServiceServer) ListForParent(context.Context, *ListLooncanForParentRequest) (*ListLooncanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListForParent not implemented")
}
func (UnimplementedLooncanServiceServer) mustEmbedUnimplementedLooncanServiceServer() {}

// UnsafeLooncanServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LooncanServiceServer will
// result in compilation errors.
type UnsafeLooncanServiceServer interface {
	mustEmbedUnimplementedLooncanServiceServer()
}

func RegisterLooncanServiceServer(s grpc.ServiceRegistrar, srv LooncanServiceServer) {
	s.RegisterService(&LooncanService_ServiceDesc, srv)
}

func _LooncanService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LooncanServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LooncanService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LooncanServiceServer).List(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LooncanService_ListForParent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLooncanForParentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LooncanServiceServer).ListForParent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LooncanService_ListForParent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LooncanServiceServer).ListForParent(ctx, req.(*ListLooncanForParentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LooncanService_ServiceDesc is the grpc.ServiceDesc for LooncanService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LooncanService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "microservice.LooncanService",
	HandlerType: (*LooncanServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _LooncanService_List_Handler,
		},
		{
			MethodName: "ListForParent",
			Handler:    _LooncanService_ListForParent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/looncan.proto",
}
