package grpcio

import (
	"io"

	"github.com/frantjc/sequence/api/types"
	"google.golang.org/grpc"
)

type LogStreamClient interface {
	Recv() (*types.Log, error)
	grpc.ClientStream
}

func MultiplexLogStream(stream LogStreamClient, stdout, stderr io.Writer) error {
	for {
		l, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		} else if l != nil {
			if len(l.Out) > 0 {
				stdout.Write([]byte(l.Out))
			}
			if len(l.Err) > 0 {
				stderr.Write([]byte(l.Err))
			}
		}
	}
}
