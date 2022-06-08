// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: runtime/v1/container.proto

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
	// ContainerServiceName is the fully-qualified name of the ContainerService service.
	ContainerServiceName = "runtime.v1.ContainerService"
)

// ContainerServiceClient is a client for the runtime.v1.ContainerService service.
type ContainerServiceClient interface {
	CreateContainer(context.Context, *connect_go.Request[v1.CreateContainerRequest]) (*connect_go.Response[v1.CreateContainerResponse], error)
	GetContainer(context.Context, *connect_go.Request[v1.GetContainerRequest]) (*connect_go.Response[v1.GetContainerResponse], error)
	ExecContainer(context.Context, *connect_go.Request[v1.ExecContainerRequest]) (*connect_go.ServerStreamForClient[v1.ExecContainerResponse], error)
	StartContainer(context.Context, *connect_go.Request[v1.StartContainerRequest]) (*connect_go.Response[v1.StartContainerResponse], error)
	AttachContainer(context.Context, *connect_go.Request[v1.AttachContainerRequest]) (*connect_go.ServerStreamForClient[v1.AttachContainerResponse], error)
	RemoveContainer(context.Context, *connect_go.Request[v1.RemoveContainerRequest]) (*connect_go.Response[v1.RemoveContainerResponse], error)
	PruneContainers(context.Context, *connect_go.Request[v1.PruneContainersRequest]) (*connect_go.Response[v1.PruneContainersResponse], error)
	CopyToContainer(context.Context, *connect_go.Request[v1.CopyToContainerRequest]) (*connect_go.Response[v1.CopyToContainerResponse], error)
	CopyFromContainer(context.Context, *connect_go.Request[v1.CopyFromContainerRequest]) (*connect_go.Response[v1.CopyFromContainerResponse], error)
}

// NewContainerServiceClient constructs a client for the runtime.v1.ContainerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewContainerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ContainerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &containerServiceClient{
		createContainer: connect_go.NewClient[v1.CreateContainerRequest, v1.CreateContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/CreateContainer",
			opts...,
		),
		getContainer: connect_go.NewClient[v1.GetContainerRequest, v1.GetContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/GetContainer",
			opts...,
		),
		execContainer: connect_go.NewClient[v1.ExecContainerRequest, v1.ExecContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/ExecContainer",
			opts...,
		),
		startContainer: connect_go.NewClient[v1.StartContainerRequest, v1.StartContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/StartContainer",
			opts...,
		),
		attachContainer: connect_go.NewClient[v1.AttachContainerRequest, v1.AttachContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/AttachContainer",
			opts...,
		),
		removeContainer: connect_go.NewClient[v1.RemoveContainerRequest, v1.RemoveContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/RemoveContainer",
			opts...,
		),
		pruneContainers: connect_go.NewClient[v1.PruneContainersRequest, v1.PruneContainersResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/PruneContainers",
			opts...,
		),
		copyToContainer: connect_go.NewClient[v1.CopyToContainerRequest, v1.CopyToContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/CopyToContainer",
			opts...,
		),
		copyFromContainer: connect_go.NewClient[v1.CopyFromContainerRequest, v1.CopyFromContainerResponse](
			httpClient,
			baseURL+"/runtime.v1.ContainerService/CopyFromContainer",
			opts...,
		),
	}
}

// containerServiceClient implements ContainerServiceClient.
type containerServiceClient struct {
	createContainer   *connect_go.Client[v1.CreateContainerRequest, v1.CreateContainerResponse]
	getContainer      *connect_go.Client[v1.GetContainerRequest, v1.GetContainerResponse]
	execContainer     *connect_go.Client[v1.ExecContainerRequest, v1.ExecContainerResponse]
	startContainer    *connect_go.Client[v1.StartContainerRequest, v1.StartContainerResponse]
	attachContainer   *connect_go.Client[v1.AttachContainerRequest, v1.AttachContainerResponse]
	removeContainer   *connect_go.Client[v1.RemoveContainerRequest, v1.RemoveContainerResponse]
	pruneContainers   *connect_go.Client[v1.PruneContainersRequest, v1.PruneContainersResponse]
	copyToContainer   *connect_go.Client[v1.CopyToContainerRequest, v1.CopyToContainerResponse]
	copyFromContainer *connect_go.Client[v1.CopyFromContainerRequest, v1.CopyFromContainerResponse]
}

// CreateContainer calls runtime.v1.ContainerService.CreateContainer.
func (c *containerServiceClient) CreateContainer(ctx context.Context, req *connect_go.Request[v1.CreateContainerRequest]) (*connect_go.Response[v1.CreateContainerResponse], error) {
	return c.createContainer.CallUnary(ctx, req)
}

// GetContainer calls runtime.v1.ContainerService.GetContainer.
func (c *containerServiceClient) GetContainer(ctx context.Context, req *connect_go.Request[v1.GetContainerRequest]) (*connect_go.Response[v1.GetContainerResponse], error) {
	return c.getContainer.CallUnary(ctx, req)
}

// ExecContainer calls runtime.v1.ContainerService.ExecContainer.
func (c *containerServiceClient) ExecContainer(ctx context.Context, req *connect_go.Request[v1.ExecContainerRequest]) (*connect_go.ServerStreamForClient[v1.ExecContainerResponse], error) {
	return c.execContainer.CallServerStream(ctx, req)
}

// StartContainer calls runtime.v1.ContainerService.StartContainer.
func (c *containerServiceClient) StartContainer(ctx context.Context, req *connect_go.Request[v1.StartContainerRequest]) (*connect_go.Response[v1.StartContainerResponse], error) {
	return c.startContainer.CallUnary(ctx, req)
}

// AttachContainer calls runtime.v1.ContainerService.AttachContainer.
func (c *containerServiceClient) AttachContainer(ctx context.Context, req *connect_go.Request[v1.AttachContainerRequest]) (*connect_go.ServerStreamForClient[v1.AttachContainerResponse], error) {
	return c.attachContainer.CallServerStream(ctx, req)
}

// RemoveContainer calls runtime.v1.ContainerService.RemoveContainer.
func (c *containerServiceClient) RemoveContainer(ctx context.Context, req *connect_go.Request[v1.RemoveContainerRequest]) (*connect_go.Response[v1.RemoveContainerResponse], error) {
	return c.removeContainer.CallUnary(ctx, req)
}

// PruneContainers calls runtime.v1.ContainerService.PruneContainers.
func (c *containerServiceClient) PruneContainers(ctx context.Context, req *connect_go.Request[v1.PruneContainersRequest]) (*connect_go.Response[v1.PruneContainersResponse], error) {
	return c.pruneContainers.CallUnary(ctx, req)
}

// CopyToContainer calls runtime.v1.ContainerService.CopyToContainer.
func (c *containerServiceClient) CopyToContainer(ctx context.Context, req *connect_go.Request[v1.CopyToContainerRequest]) (*connect_go.Response[v1.CopyToContainerResponse], error) {
	return c.copyToContainer.CallUnary(ctx, req)
}

// CopyFromContainer calls runtime.v1.ContainerService.CopyFromContainer.
func (c *containerServiceClient) CopyFromContainer(ctx context.Context, req *connect_go.Request[v1.CopyFromContainerRequest]) (*connect_go.Response[v1.CopyFromContainerResponse], error) {
	return c.copyFromContainer.CallUnary(ctx, req)
}

// ContainerServiceHandler is an implementation of the runtime.v1.ContainerService service.
type ContainerServiceHandler interface {
	CreateContainer(context.Context, *connect_go.Request[v1.CreateContainerRequest]) (*connect_go.Response[v1.CreateContainerResponse], error)
	GetContainer(context.Context, *connect_go.Request[v1.GetContainerRequest]) (*connect_go.Response[v1.GetContainerResponse], error)
	ExecContainer(context.Context, *connect_go.Request[v1.ExecContainerRequest], *connect_go.ServerStream[v1.ExecContainerResponse]) error
	StartContainer(context.Context, *connect_go.Request[v1.StartContainerRequest]) (*connect_go.Response[v1.StartContainerResponse], error)
	AttachContainer(context.Context, *connect_go.Request[v1.AttachContainerRequest], *connect_go.ServerStream[v1.AttachContainerResponse]) error
	RemoveContainer(context.Context, *connect_go.Request[v1.RemoveContainerRequest]) (*connect_go.Response[v1.RemoveContainerResponse], error)
	PruneContainers(context.Context, *connect_go.Request[v1.PruneContainersRequest]) (*connect_go.Response[v1.PruneContainersResponse], error)
	CopyToContainer(context.Context, *connect_go.Request[v1.CopyToContainerRequest]) (*connect_go.Response[v1.CopyToContainerResponse], error)
	CopyFromContainer(context.Context, *connect_go.Request[v1.CopyFromContainerRequest]) (*connect_go.Response[v1.CopyFromContainerResponse], error)
}

// NewContainerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewContainerServiceHandler(svc ContainerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/runtime.v1.ContainerService/CreateContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/CreateContainer",
		svc.CreateContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/GetContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/GetContainer",
		svc.GetContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/ExecContainer", connect_go.NewServerStreamHandler(
		"/runtime.v1.ContainerService/ExecContainer",
		svc.ExecContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/StartContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/StartContainer",
		svc.StartContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/AttachContainer", connect_go.NewServerStreamHandler(
		"/runtime.v1.ContainerService/AttachContainer",
		svc.AttachContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/RemoveContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/RemoveContainer",
		svc.RemoveContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/PruneContainers", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/PruneContainers",
		svc.PruneContainers,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/CopyToContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/CopyToContainer",
		svc.CopyToContainer,
		opts...,
	))
	mux.Handle("/runtime.v1.ContainerService/CopyFromContainer", connect_go.NewUnaryHandler(
		"/runtime.v1.ContainerService/CopyFromContainer",
		svc.CopyFromContainer,
		opts...,
	))
	return "/runtime.v1.ContainerService/", mux
}

// UnimplementedContainerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedContainerServiceHandler struct{}

func (UnimplementedContainerServiceHandler) CreateContainer(context.Context, *connect_go.Request[v1.CreateContainerRequest]) (*connect_go.Response[v1.CreateContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.CreateContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) GetContainer(context.Context, *connect_go.Request[v1.GetContainerRequest]) (*connect_go.Response[v1.GetContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.GetContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) ExecContainer(context.Context, *connect_go.Request[v1.ExecContainerRequest], *connect_go.ServerStream[v1.ExecContainerResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.ExecContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) StartContainer(context.Context, *connect_go.Request[v1.StartContainerRequest]) (*connect_go.Response[v1.StartContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.StartContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) AttachContainer(context.Context, *connect_go.Request[v1.AttachContainerRequest], *connect_go.ServerStream[v1.AttachContainerResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.AttachContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) RemoveContainer(context.Context, *connect_go.Request[v1.RemoveContainerRequest]) (*connect_go.Response[v1.RemoveContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.RemoveContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) PruneContainers(context.Context, *connect_go.Request[v1.PruneContainersRequest]) (*connect_go.Response[v1.PruneContainersResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.PruneContainers is not implemented"))
}

func (UnimplementedContainerServiceHandler) CopyToContainer(context.Context, *connect_go.Request[v1.CopyToContainerRequest]) (*connect_go.Response[v1.CopyToContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.CopyToContainer is not implemented"))
}

func (UnimplementedContainerServiceHandler) CopyFromContainer(context.Context, *connect_go.Request[v1.CopyFromContainerRequest]) (*connect_go.Response[v1.CopyFromContainerResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("runtime.v1.ContainerService.CopyFromContainer is not implemented"))
}
