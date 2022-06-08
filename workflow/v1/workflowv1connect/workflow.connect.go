// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflow/v1/workflow.proto

package workflowv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/frantjc/sequence/workflow/v1"
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
	// WorkflowServiceName is the fully-qualified name of the WorkflowService service.
	WorkflowServiceName = "workflow.v1.WorkflowService"
)

// WorkflowServiceClient is a client for the workflow.v1.WorkflowService service.
type WorkflowServiceClient interface {
	RunWorkflow(context.Context, *connect_go.Request[v1.RunWorkflowRequest]) (*connect_go.ServerStreamForClient[v1.RunWorkflowResponse], error)
	RunJob(context.Context, *connect_go.Request[v1.RunJobRequest]) (*connect_go.ServerStreamForClient[v1.RunJobResponse], error)
	RunStep(context.Context, *connect_go.Request[v1.RunStepRequest]) (*connect_go.ServerStreamForClient[v1.RunStepResponse], error)
}

// NewWorkflowServiceClient constructs a client for the workflow.v1.WorkflowService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewWorkflowServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) WorkflowServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &workflowServiceClient{
		runWorkflow: connect_go.NewClient[v1.RunWorkflowRequest, v1.RunWorkflowResponse](
			httpClient,
			baseURL+"/workflow.v1.WorkflowService/RunWorkflow",
			opts...,
		),
		runJob: connect_go.NewClient[v1.RunJobRequest, v1.RunJobResponse](
			httpClient,
			baseURL+"/workflow.v1.WorkflowService/RunJob",
			opts...,
		),
		runStep: connect_go.NewClient[v1.RunStepRequest, v1.RunStepResponse](
			httpClient,
			baseURL+"/workflow.v1.WorkflowService/RunStep",
			opts...,
		),
	}
}

// workflowServiceClient implements WorkflowServiceClient.
type workflowServiceClient struct {
	runWorkflow *connect_go.Client[v1.RunWorkflowRequest, v1.RunWorkflowResponse]
	runJob      *connect_go.Client[v1.RunJobRequest, v1.RunJobResponse]
	runStep     *connect_go.Client[v1.RunStepRequest, v1.RunStepResponse]
}

// RunWorkflow calls workflow.v1.WorkflowService.RunWorkflow.
func (c *workflowServiceClient) RunWorkflow(ctx context.Context, req *connect_go.Request[v1.RunWorkflowRequest]) (*connect_go.ServerStreamForClient[v1.RunWorkflowResponse], error) {
	return c.runWorkflow.CallServerStream(ctx, req)
}

// RunJob calls workflow.v1.WorkflowService.RunJob.
func (c *workflowServiceClient) RunJob(ctx context.Context, req *connect_go.Request[v1.RunJobRequest]) (*connect_go.ServerStreamForClient[v1.RunJobResponse], error) {
	return c.runJob.CallServerStream(ctx, req)
}

// RunStep calls workflow.v1.WorkflowService.RunStep.
func (c *workflowServiceClient) RunStep(ctx context.Context, req *connect_go.Request[v1.RunStepRequest]) (*connect_go.ServerStreamForClient[v1.RunStepResponse], error) {
	return c.runStep.CallServerStream(ctx, req)
}

// WorkflowServiceHandler is an implementation of the workflow.v1.WorkflowService service.
type WorkflowServiceHandler interface {
	RunWorkflow(context.Context, *connect_go.Request[v1.RunWorkflowRequest], *connect_go.ServerStream[v1.RunWorkflowResponse]) error
	RunJob(context.Context, *connect_go.Request[v1.RunJobRequest], *connect_go.ServerStream[v1.RunJobResponse]) error
	RunStep(context.Context, *connect_go.Request[v1.RunStepRequest], *connect_go.ServerStream[v1.RunStepResponse]) error
}

// NewWorkflowServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewWorkflowServiceHandler(svc WorkflowServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/workflow.v1.WorkflowService/RunWorkflow", connect_go.NewServerStreamHandler(
		"/workflow.v1.WorkflowService/RunWorkflow",
		svc.RunWorkflow,
		opts...,
	))
	mux.Handle("/workflow.v1.WorkflowService/RunJob", connect_go.NewServerStreamHandler(
		"/workflow.v1.WorkflowService/RunJob",
		svc.RunJob,
		opts...,
	))
	mux.Handle("/workflow.v1.WorkflowService/RunStep", connect_go.NewServerStreamHandler(
		"/workflow.v1.WorkflowService/RunStep",
		svc.RunStep,
		opts...,
	))
	return "/workflow.v1.WorkflowService/", mux
}

// UnimplementedWorkflowServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedWorkflowServiceHandler struct{}

func (UnimplementedWorkflowServiceHandler) RunWorkflow(context.Context, *connect_go.Request[v1.RunWorkflowRequest], *connect_go.ServerStream[v1.RunWorkflowResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("workflow.v1.WorkflowService.RunWorkflow is not implemented"))
}

func (UnimplementedWorkflowServiceHandler) RunJob(context.Context, *connect_go.Request[v1.RunJobRequest], *connect_go.ServerStream[v1.RunJobResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("workflow.v1.WorkflowService.RunJob is not implemented"))
}

func (UnimplementedWorkflowServiceHandler) RunStep(context.Context, *connect_go.Request[v1.RunStepRequest], *connect_go.ServerStream[v1.RunStepResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("workflow.v1.WorkflowService.RunStep is not implemented"))
}
