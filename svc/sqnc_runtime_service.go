package svc

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/internal/rpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
)

type SqncRuntimeServiceHandler struct {
	sqnc.UnimplementedRuntimeServiceHandler
	Runtime runtime.Runtime
}

func (h *SqncRuntimeServiceHandler) CreateContainer(ctx context.Context, req *connect.Request[sqnc.CreateContainerRequest]) (*connect.Response[sqnc.CreateContainerResponse], error) {
	container, err := h.Runtime.CreateContainer(ctx, req.Msg.GetSpec())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.CreateContainerResponse{
		Container: &sqnc.Container{
			Id: container.GetID(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) GetContainer(ctx context.Context, req *connect.Request[sqnc.GetContainerRequest]) (*connect.Response[sqnc.GetContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&sqnc.GetContainerResponse{
		Container: &sqnc.Container{
			Id: container.GetID(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) ExecContainer(ctx context.Context, req *connect.Request[sqnc.ExecContainerRequest], stream *connect.ServerStream[sqnc.ExecContainerResponse]) error {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.Exec(ctx, &runtime.Exec{
		Cmd: req.Msg.GetExec().GetCmd(),
	}, &runtime.Streams{
		In: os.Stdin,
		Out: rpcio.NewServerStreamWriter[sqnc.ExecContainerResponse](stream, func(b []byte) *sqnc.ExecContainerResponse {
			return &sqnc.ExecContainerResponse{
				Log: &rpcio.Log{
					Data:   b,
					Stream: int32(rpcio.StreamStdout),
				},
			}
		}),
		Err: rpcio.NewServerStreamWriter[sqnc.ExecContainerResponse](stream, func(b []byte) *sqnc.ExecContainerResponse {
			return &sqnc.ExecContainerResponse{
				Log: &rpcio.Log{
					Data:   b,
					Stream: int32(rpcio.StreamStderr),
				},
			}
		}),
	}); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}

func (h *SqncRuntimeServiceHandler) StartContainer(ctx context.Context, req *connect.Request[sqnc.StartContainerRequest]) (*connect.Response[sqnc.StartContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.Start(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.StartContainerResponse{
		Container: &sqnc.Container{
			Id: container.GetID(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) AttachContainer(ctx context.Context, req *connect.Request[sqnc.AttachContainerRequest], stream *connect.ServerStream[sqnc.AttachContainerResponse]) error {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.Attach(ctx, runtime.NewStreams(
		os.Stdin,
		rpcio.NewServerStreamWriter[sqnc.AttachContainerResponse](stream, func(b []byte) *sqnc.AttachContainerResponse {
			return &sqnc.AttachContainerResponse{
				Log: &rpcio.Log{
					Data:   b,
					Stream: int32(rpcio.StreamStdout),
				},
			}
		}),
		rpcio.NewServerStreamWriter[sqnc.AttachContainerResponse](stream, func(b []byte) *sqnc.AttachContainerResponse {
			return &sqnc.AttachContainerResponse{
				Log: &rpcio.Log{
					Data:   b,
					Stream: int32(rpcio.StreamStderr),
				},
			}
		}),
	)); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}

func (h *SqncRuntimeServiceHandler) StopContainer(ctx context.Context, req *connect.Request[sqnc.StopContainerRequest]) (*connect.Response[sqnc.StopContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.Stop(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.StopContainerResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) RemoveContainer(ctx context.Context, req *connect.Request[sqnc.RemoveContainerRequest]) (*connect.Response[sqnc.RemoveContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.Remove(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.RemoveContainerResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) PruneContainers(ctx context.Context, _ *connect.Request[sqnc.PruneContainersRequest]) (*connect.Response[sqnc.PruneContainersResponse], error) {
	if err := h.Runtime.PruneContainers(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.PruneContainersResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) CopyToContainer(ctx context.Context, req *connect.Request[sqnc.CopyToContainerRequest]) (*connect.Response[sqnc.CopyToContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := container.CopyTo(ctx, bytes.NewReader(req.Msg.GetContent()), req.Msg.GetDestination()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.CopyToContainerResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) CopyFromContainer(ctx context.Context, req *connect.Request[sqnc.CopyFromContainerRequest]) (*connect.Response[sqnc.CopyFromContainerResponse], error) {
	container, err := h.Runtime.GetContainer(ctx, req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	r, err := container.CopyFrom(ctx, req.Msg.GetSource())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.CopyFromContainerResponse{
		Content: content,
	}), nil
}

func (h *SqncRuntimeServiceHandler) PullImage(ctx context.Context, req *connect.Request[sqnc.PullImageRequest]) (*connect.Response[sqnc.PullImageResponse], error) {
	image, err := h.Runtime.PullImage(ctx, req.Msg.GetRef())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.PullImageResponse{
		Image: &sqnc.Image{
			Ref: image.GetRef(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) PruneImages(ctx context.Context, req *connect.Request[sqnc.PruneImagesRequest]) (*connect.Response[sqnc.PruneImagesResponse], error) {
	if err := h.Runtime.PruneImages(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.PruneImagesResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) CreateVolume(ctx context.Context, req *connect.Request[sqnc.CreateVolumeRequest]) (*connect.Response[sqnc.CreateVolumeResponse], error) {
	volume, err := h.Runtime.CreateVolume(ctx, req.Msg.GetName())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.CreateVolumeResponse{
		Volume: &sqnc.Volume{
			Source: volume.GetSource(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) GetVolume(ctx context.Context, req *connect.Request[sqnc.GetVolumeRequest]) (*connect.Response[sqnc.GetVolumeResponse], error) {
	volume, err := h.Runtime.GetVolume(ctx, req.Msg.GetName())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&sqnc.GetVolumeResponse{
		Volume: &sqnc.Volume{
			Source: volume.GetSource(),
		},
	}), nil
}

func (h *SqncRuntimeServiceHandler) RemoveVolume(ctx context.Context, req *connect.Request[sqnc.RemoveVolumeRequest]) (*connect.Response[sqnc.RemoveVolumeResponse], error) {
	volume, err := h.Runtime.GetVolume(ctx, req.Msg.GetName())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err = volume.Remove(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.RemoveVolumeResponse{}), nil
}

func (h *SqncRuntimeServiceHandler) PruneVolumes(ctx context.Context, _ *connect.Request[sqnc.PruneVolumesRequest]) (*connect.Response[sqnc.PruneVolumesResponse], error) {
	if err := h.Runtime.PruneVolumes(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&sqnc.PruneVolumesResponse{}), nil
}
