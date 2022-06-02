package grpcio

import "io"

func NewLogStreamReader(stream LogStreamClient) io.Reader {
	return &logStreamReader{stream}
}

type logStreamReader struct {
	s LogStreamClient
}

func (r *logStreamReader) Read(p []byte) (int, error) {
	log, err := r.s.Recv()
	if err != nil {
		return 0, err
	}

	return copy(p, log.Data), nil
}

func (r *logStreamReader) Close() error {
	return r.s.CloseSend()
}

var _ io.ReadCloser = &logStreamReader{}
