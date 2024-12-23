// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: estimation/v1/estimation.proto

package estimationv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	EstimationService_EstimateUsers_FullMethodName = "/estimation.v1.EstimationService/EstimateUsers"
)

// EstimationServiceClient is the client API for EstimationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EstimationServiceClient interface {
	EstimateUsers(ctx context.Context, in *EstimateUsersRequest, opts ...grpc.CallOption) (*EstimateUsersResponse, error)
}

type estimationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEstimationServiceClient(cc grpc.ClientConnInterface) EstimationServiceClient {
	return &estimationServiceClient{cc}
}

func (c *estimationServiceClient) EstimateUsers(ctx context.Context, in *EstimateUsersRequest, opts ...grpc.CallOption) (*EstimateUsersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EstimateUsersResponse)
	err := c.cc.Invoke(ctx, EstimationService_EstimateUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EstimationServiceServer is the server API for EstimationService service.
// All implementations must embed UnimplementedEstimationServiceServer
// for forward compatibility
type EstimationServiceServer interface {
	EstimateUsers(context.Context, *EstimateUsersRequest) (*EstimateUsersResponse, error)
	mustEmbedUnimplementedEstimationServiceServer()
}

// UnimplementedEstimationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEstimationServiceServer struct {
}

func (UnimplementedEstimationServiceServer) EstimateUsers(context.Context, *EstimateUsersRequest) (*EstimateUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimateUsers not implemented")
}
func (UnimplementedEstimationServiceServer) mustEmbedUnimplementedEstimationServiceServer() {}

// UnsafeEstimationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EstimationServiceServer will
// result in compilation errors.
type UnsafeEstimationServiceServer interface {
	mustEmbedUnimplementedEstimationServiceServer()
}

func RegisterEstimationServiceServer(s grpc.ServiceRegistrar, srv EstimationServiceServer) {
	s.RegisterService(&EstimationService_ServiceDesc, srv)
}

func _EstimationService_EstimateUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EstimateUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstimationServiceServer).EstimateUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EstimationService_EstimateUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstimationServiceServer).EstimateUsers(ctx, req.(*EstimateUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EstimationService_ServiceDesc is the grpc.ServiceDesc for EstimationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EstimationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "estimation.v1.EstimationService",
	HandlerType: (*EstimationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EstimateUsers",
			Handler:    _EstimationService_EstimateUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "estimation/v1/estimation.proto",
}
