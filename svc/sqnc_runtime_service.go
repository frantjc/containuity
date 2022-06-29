package svc

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime/sqnc"
)

type SqncRuntimeServiceHandler struct {
	sqnc.UnimplementedRuntimeServiceHandler
}

func (*SqncRuntimeServiceHandler) CreateContainer(context.Context, *connect.Request[sqnc.CreateContainerRequest]) (*connect.Response[sqnc.CreateContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.CreateContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) GetContainer(context.Context, *connect.Request[sqnc.GetContainerRequest]) (*connect.Response[sqnc.GetContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.GetContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) ExecContainer(context.Context, *connect.Request[sqnc.ExecContainerRequest], *connect.ServerStream[sqnc.ExecContainerResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.ExecContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) StartContainer(context.Context, *connect.Request[sqnc.StartContainerRequest]) (*connect.Response[sqnc.StartContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.StartContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) AttachContainer(context.Context, *connect.Request[sqnc.AttachContainerRequest], *connect.ServerStream[sqnc.AttachContainerResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.AttachContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) RemoveContainer(context.Context, *connect.Request[sqnc.RemoveContainerRequest]) (*connect.Response[sqnc.RemoveContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.RemoveContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) PruneContainers(context.Context, *connect.Request[sqnc.PruneContainersRequest]) (*connect.Response[sqnc.PruneContainersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.PruneContainers is not implemented"))
}

func (*SqncRuntimeServiceHandler) CopyToContainer(context.Context, *connect.Request[sqnc.CopyToContainerRequest]) (*connect.Response[sqnc.CopyToContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.CopyToContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) CopyFromContainer(context.Context, *connect.Request[sqnc.CopyFromContainerRequest]) (*connect.Response[sqnc.CopyFromContainerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.CopyFromContainer is not implemented"))
}

func (*SqncRuntimeServiceHandler) PullImage(context.Context, *connect.Request[sqnc.PullImageRequest]) (*connect.Response[sqnc.PullImageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.PullImage is not implemented"))
}

func (*SqncRuntimeServiceHandler) PruneImages(context.Context, *connect.Request[sqnc.PruneImagesRequest]) (*connect.Response[sqnc.PruneImagesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.PruneImages is not implemented"))
}

func (*SqncRuntimeServiceHandler) CreateVolume(context.Context, *connect.Request[sqnc.CreateVolumeRequest]) (*connect.Response[sqnc.CreateVolumeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.CreateVolume is not implemented"))
}

func (*SqncRuntimeServiceHandler) GetVolume(context.Context, *connect.Request[sqnc.GetVolumeRequest]) (*connect.Response[sqnc.GetVolumeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.GetVolume is not implemented"))
}

func (*SqncRuntimeServiceHandler) RemoveVolume(context.Context, *connect.Request[sqnc.RemoveVolumeRequest]) (*connect.Response[sqnc.RemoveVolumeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.RemoveVolume is not implemented"))
}

func (*SqncRuntimeServiceHandler) PruneVolumes(context.Context, *connect.Request[sqnc.PruneVolumesRequest]) (*connect.Response[sqnc.PruneVolumesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sequence.runtime.sqnc.RuntimeService.PruneVolumes is not implemented"))
}
