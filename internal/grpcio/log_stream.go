package grpcio

import (
	"context"
	"io"

	"github.com/frantjc/sequence/pb/types"
	"google.golang.org/grpc/metadata"
)

type LogStream interface {
	LogStreamClient
	LogStreamServer
}

func NewLogStream(ctx context.Context) LogStream {
	return &logStream{
		ctx:    ctx,
		logC:   make(chan *types.Log),
		closed: false,
		errC:   make(chan error, 1),
	}
}

type logStream struct {
	ctx    context.Context
	logC   chan *types.Log
	closed bool
	errC   chan error
}

var _ LogStream = &logStream{}

func (s *logStream) Recv() (*types.Log, error) {
	if s.closed {
		return nil, io.EOF
	}
	select {
	case l := <-s.logC:
		return l, nil
	case err := <-s.errC:
		return nil, err
	}
}

func (s *logStream) Send(l *types.Log) error {
	if !s.closed {
		s.logC <- l
		return nil
	}
	return ErrStreamClosed
}

func (s *logStream) Context() context.Context {
	return s.ctx
}

func (s *logStream) CloseSend() error {
	if !s.closed {
		close(s.logC)
		s.closed = true
		return nil
	}
	return ErrStreamClosed
}

func (s *logStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (s *logStream) Trailer() metadata.MD {
	return nil
}

func (s *logStream) SendMsg(m interface{}) error {
	return nil
}

func (s *logStream) RecvMsg(m interface{}) error {
	return nil
}
