// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package volume

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

// VolumeClient is the client API for Volume service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VolumeClient interface {
	CreateVolume(ctx context.Context, in *CreateVolumeRequest, opts ...grpc.CallOption) (*CreateVolumeResponse, error)
	GetVolume(ctx context.Context, in *GetVolumeRequest, opts ...grpc.CallOption) (*GetVolumeResponse, error)
	RemoveVolume(ctx context.Context, in *RemoveVolumeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PruneVolumes(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type volumeClient struct {
	cc grpc.ClientConnInterface
}

func NewVolumeClient(cc grpc.ClientConnInterface) VolumeClient {
	return &volumeClient{cc}
}

func (c *volumeClient) CreateVolume(ctx context.Context, in *CreateVolumeRequest, opts ...grpc.CallOption) (*CreateVolumeResponse, error) {
	out := new(CreateVolumeResponse)
	err := c.cc.Invoke(ctx, "/sequence.v1.volume.Volume/CreateVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) GetVolume(ctx context.Context, in *GetVolumeRequest, opts ...grpc.CallOption) (*GetVolumeResponse, error) {
	out := new(GetVolumeResponse)
	err := c.cc.Invoke(ctx, "/sequence.v1.volume.Volume/GetVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) RemoveVolume(ctx context.Context, in *RemoveVolumeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sequence.v1.volume.Volume/RemoveVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volumeClient) PruneVolumes(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sequence.v1.volume.Volume/PruneVolumes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VolumeServer is the server API for Volume service.
// All implementations must embed UnimplementedVolumeServer
// for forward compatibility
type VolumeServer interface {
	CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)
	GetVolume(context.Context, *GetVolumeRequest) (*GetVolumeResponse, error)
	RemoveVolume(context.Context, *RemoveVolumeRequest) (*emptypb.Empty, error)
	PruneVolumes(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedVolumeServer()
}

// UnimplementedVolumeServer must be embedded to have forward compatible implementations.
type UnimplementedVolumeServer struct {
}

func (UnimplementedVolumeServer) CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVolume not implemented")
}
func (UnimplementedVolumeServer) GetVolume(context.Context, *GetVolumeRequest) (*GetVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVolume not implemented")
}
func (UnimplementedVolumeServer) RemoveVolume(context.Context, *RemoveVolumeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveVolume not implemented")
}
func (UnimplementedVolumeServer) PruneVolumes(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PruneVolumes not implemented")
}
func (UnimplementedVolumeServer) mustEmbedUnimplementedVolumeServer() {}

// UnsafeVolumeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VolumeServer will
// result in compilation errors.
type UnsafeVolumeServer interface {
	mustEmbedUnimplementedVolumeServer()
}

func RegisterVolumeServer(s grpc.ServiceRegistrar, srv VolumeServer) {
	s.RegisterService(&Volume_ServiceDesc, srv)
}

func _Volume_CreateVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).CreateVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.volume.Volume/CreateVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).CreateVolume(ctx, req.(*CreateVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_GetVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).GetVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.volume.Volume/GetVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).GetVolume(ctx, req.(*GetVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_RemoveVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).RemoveVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.volume.Volume/RemoveVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).RemoveVolume(ctx, req.(*RemoveVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Volume_PruneVolumes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolumeServer).PruneVolumes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequence.v1.volume.Volume/PruneVolumes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolumeServer).PruneVolumes(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Volume_ServiceDesc is the grpc.ServiceDesc for Volume service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Volume_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sequence.v1.volume.Volume",
	HandlerType: (*VolumeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateVolume",
			Handler:    _Volume_CreateVolume_Handler,
		},
		{
			MethodName: "GetVolume",
			Handler:    _Volume_GetVolume_Handler,
		},
		{
			MethodName: "RemoveVolume",
			Handler:    _Volume_RemoveVolume_Handler,
		},
		{
			MethodName: "PruneVolumes",
			Handler:    _Volume_PruneVolumes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/volume/volume.proto",
}
