// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: runtime/v1/volume.proto

package runtimev1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/frantjc/sequence/runtime/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// VolumeServiceName is the fully-qualified name of the VolumeService service.
	VolumeServiceName = "runtime.v1.VolumeService"
)

// VolumeServiceClient is a client for the runtime.v1.VolumeService service.
type VolumeServiceClient interface {
	CreateVolume(context.Context, *connect_go.Request[v1.CreateVolumeRequest]) (*connect_go.Response[v1.CreateVolumeResponse], error)
	GetVolume(context.Context, *connect_go.Request[v1.GetVolumeRequest]) (*connect_go.Response[v1.GetVolumeResponse], error)
	RemoveVolume(context.Context, *connect_go.Request[v1.RemoveVolumeRequest]) (*connect_go.Response[v1.RemoveVolumeResponse], error)
	PruneVolumes(context.Context, *connect_go.Request[v1.PruneVolumesRequest]) (*connect_go.Response[v1.PruneVolumesResponse], error)
}

// NewVolumeServiceClient constructs a client for the runtime.v1.VolumeService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewVolumeServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) VolumeServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &volumeServiceClient{
		createVolume: connect_go.NewClient[v1.CreateVolumeRequest, v1.CreateVolumeResponse](
			httpClient,
			baseURL+"/runtime.v1.VolumeService/CreateVolume",
			opts...,
		),
		getVolume: connect_go.NewClient[v1.GetVolumeRequest, v1.GetVolumeResponse](
			httpClient,
			baseURL+"/runtime.v1.VolumeService/GetVolume",
			opts...,
		),
		removeVolume: connect_go.NewClient[v1.RemoveVolumeRequest, v1.RemoveVolumeResponse](
			httpClient,
			baseURL+"/runtime.v1.VolumeService/RemoveVolume",
			opts...,
		),
		pruneVolumes: connect_go.NewClient[v1.PruneVolumesRequest, v1.PruneVolumesResponse](
			httpClient,
			baseURL+"/runtime.v1.VolumeService/PruneVolumes",
			opts...,
		),
	}
}

// volumeServiceClient implements VolumeServiceClient.
type volumeServiceClient struct {
	createVolume *connect_go.Client[v1.CreateVolumeRequest, v1.CreateVolumeResponse]
	getVolume    *connect_go.Client[v1.GetVolumeRequest, v1.GetVolumeResponse]
	removeVolume *connect_go.Client[v1.RemoveVolumeRequest, v1.RemoveVolumeResponse]
	pruneVolumes *connect_go.Client[v1.PruneVolumesRequest, v1.PruneVolumesResponse]
}

// CreateVolume calls runtime.v1.VolumeService.CreateVolume.
func (c *volumeServiceClient) CreateVolume(ctx context.Context, req *connect_go.Request[v1.CreateVolumeRequest]) (*connect_go.Response[v1.CreateVolumeResponse], error) {
	return c.createVolume.CallUnary(ctx, req)
}

// GetVolume calls runtime.v1.VolumeService.GetVolume.
func (c *volumeServiceClient) GetVolume(ctx context.Context, req *connect_go.Request[v1.GetVolumeRequest]) (*connect_go.Response[v1.GetVolumeResponse], error) {
	return c.getVolume.CallUnary(ctx, req)
}

// RemoveVolume calls runtime.v1.VolumeService.RemoveVolume.
func (c *volumeServiceClient) RemoveVolume(ctx context.Context, req *connect_go.Request[v1.RemoveVolumeRequest]) (*connect_go.Response[v1.RemoveVolumeResponse], error) {
	return c.removeVolume.CallUnary(ctx, req)
}

// PruneVolumes calls runtime.v1.VolumeService.PruneVolumes.
func (c *volumeServiceClient) PruneVolumes(ctx context.Context, req *connect_go.Request[v1.PruneVolumesRequest]) (*connect_go.Response[v1.PruneVolumesResponse], error) {
	return c.pruneVolumes.CallUnary(ctx, req)
}

// VolumeServiceHandler is an implementation of the runtime.v1.VolumeService service.
type VolumeServiceHandler interface {
	CreateVolume(context.Context, *connect_go.Request[v1.CreateVolumeRequest]) (*connect_go.Response[v1.CreateVolumeResponse], error)
	GetVolume(context.Context, *connect_go.Request[v1.GetVolumeRequest]) (*connect_go.Response[v1.GetVolumeResponse], error)
	RemoveVolume(context.Context, *connect_go.Request[v1.RemoveVolumeRequest]) (*connect_go.Response[v1.RemoveVolumeResponse], error)
	PruneVolumes(context.Context, *connect_go.Request[v1.PruneVolumesRequest]) (*connect_go.Response[v1.PruneVolumesResponse], error)
}

// NewVolumeServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewVolumeServiceHandler(svc VolumeServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/runtime.v1.VolumeService/CreateVolume", connect_go.NewUnaryHandler(
		"/runtime.v1.VolumeService/CreateVolume",
		svc.CreateVolume,
		opts...,
	))
	mux.Handle("/runtime.v1.VolumeService/GetVolume", connect_go.NewUnaryHandler(
		"/runtime.v1.VolumeService/GetVolume",
		svc.GetVolume,
		opts...,
	))
	mux.Handle("/runtime.v1.VolumeService/RemoveVolume", connect_go.NewUnaryHandler(
		"/runtime.v1.VolumeService/RemoveVolume",
		svc.RemoveVolume,
		opts...,
	))
	mux.Handle("/runtime.v1.VolumeService/PruneVolumes", connect_go.NewUnaryHandler(
		"/runtime.v1.VolumeService/PruneVolumes",
		svc.PruneVolumes,
		opts...,
	))
	return "/runtime.v1.VolumeService/", mux
}

// UnimplementedVolumeServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedVolumeServiceHandler struct{}

func (UnimplementedVolumeServiceHandler) CreateVolume(context.Context, *connect_go.Request[v1.CreateVolumeRequest]) (*connect_go.Response[v1.CreateVolumeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.VolumeService.CreateVolume is not implemented"))
}

func (UnimplementedVolumeServiceHandler) GetVolume(context.Context, *connect_go.Request[v1.GetVolumeRequest]) (*connect_go.Response[v1.GetVolumeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.VolumeService.GetVolume is not implemented"))
}

func (UnimplementedVolumeServiceHandler) RemoveVolume(context.Context, *connect_go.Request[v1.RemoveVolumeRequest]) (*connect_go.Response[v1.RemoveVolumeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.VolumeService.RemoveVolume is not implemented"))
}

func (UnimplementedVolumeServiceHandler) PruneVolumes(context.Context, *connect_go.Request[v1.PruneVolumesRequest]) (*connect_go.Response[v1.PruneVolumesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.VolumeService.PruneVolumes is not implemented"))
}
