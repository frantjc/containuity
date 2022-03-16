// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package container

import (
	context "context"
	types "github.com/frantjc/sequence/api/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ContainerClient is the client API for Container service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContainerClient interface {
	CreateContainer(ctx context.Context, in *CreateContainerRequest, opts ...grpc.CallOption) (*CreateContainerResponse, error)
	GetContainer(ctx context.Context, in *GetContainerRequest, opts ...grpc.CallOption) (*GetContainerResponse, error)
	ExecContainer(ctx context.Context, in *ExecContainerRequest, opts ...grpc.CallOption) (Container_ExecContainerClient, error)
}

type containerClient struct {
	cc grpc.ClientConnInterface
}

func NewContainerClient(cc grpc.ClientConnInterface) ContainerClient {
	return &containerClient{cc}
}

func (c *containerClient) CreateContainer(ctx context.Context, in *CreateContainerRequest, opts ...grpc.CallOption) (*CreateContainerResponse, error) {
	out := new(CreateContainerResponse)
	err := c.cc.Invoke(ctx, "/sequence.v1.container.Container/CreateContainer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *containerClient) GetContainer(ctx context.Context, in *GetContainerRequest, opts ...grpc.CallOption) (*GetContainerResponse, error) {
	out := new(GetContainerResponse)
	err := c.cc.Invoke(ctx, "/sequence.v1.container.Container/GetContainer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *containerClient) ExecContainer(ctx context.Context, in *ExecContainerRequest, opts ...grpc.CallOption) (Container_ExecContainerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Container_ServiceDesc.Streams[0], "/sequence.v1.container.Container/ExecContainer", opts...)
	if err != nil {
		return nil, err
	}
	x := &containerExecContainerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Container_ExecContainerClient interface {
	Recv() (*types.Log, error)
	grpc.ClientStream
}

type containerExecContainerClient struct {
	grpc.ClientStream
}

func (x *containerExecContainerClient) Recv() (*types.Log, error) {
	m := new(types.Log)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ContainerServer is the server API for Container service.
// All implementations must embed UnimplementedContainerServer
// for forward compatibility
type ContainerServer interface {
	CreateContainer(context.Context, *CreateContainerRequest) (*CreateContainerResponse, error)
	GetContainer(context.Context, *GetContainerRequest) (*GetContainerResponse, error)
	ExecContainer(*ExecContainerRequest, Container_ExecContainerServer) error
	mustEmbedUnimplementedContainerServer()
}

// UnimplementedContainerServer must be embedded to have forward compatible implementations.
type UnimplementedContainerServer struct {
}

func (UnimplementedContainerServer) CreateContainer(context.Context, *CreateContainerRequest) (*CreateContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContainer not implemented")
}
func (UnimplementedContainerServer) GetContainer(context.Context, *GetContainerRequest) (*GetContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContainer not implemented")
}
func (UnimplementedContainerServer) ExecContainer(*ExecContainerRequest, Container_ExecContainerServer) error {
	return status.Errorf(codes.Unimplemented, "method ExecContainer not implemented")
}
func (UnimplementedContainerServer) mustEmbedUnimplementedContainerServer() {}

// UnsafeContainerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContainerServer will
// result in compilation errors.
type UnsafeContainerServer interface {
	mustEmbedUnimplementedContainerServer()
}

func RegisterContainerServer(s grpc.ServiceRegistrar, srv ContainerServer) {
	s.RegisterService(&Container_ServiceDesc, srv)
}

func _Container_CreateContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContainerServer).CreateContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.container.Container/CreateContainer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContainerServer).CreateContainer(ctx, req.(*CreateContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Container_GetContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContainerServer).GetContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.container.Container/GetContainer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContainerServer).GetContainer(ctx, req.(*GetContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Container_ExecContainer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecContainerRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ContainerServer).ExecContainer(m, &containerExecContainerServer{stream})
}

type Container_ExecContainerServer interface {
	Send(*types.Log) error
	grpc.ServerStream
}

type containerExecContainerServer struct {
	grpc.ServerStream
}

func (x *containerExecContainerServer) Send(m *types.Log) error {
	return x.ServerStream.SendMsg(m)
}

// Container_ServiceDesc is the grpc.ServiceDesc for Container service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Container_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sequence.v1.container.Container",
	HandlerType: (*ContainerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateContainer",
			Handler:    _Container_CreateContainer_Handler,
		},
		{
			MethodName: "GetContainer",
			Handler:    _Container_GetContainer_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecContainer",
			Handler:       _Container_ExecContainer_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/container/container.proto",
}
