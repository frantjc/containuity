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

func DemultiplexLogStream(stream LogStreamClient, stdout, stderr io.Writer) error {
	for {
		l, err := stream.Recv()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case l != nil:
			if len(l.Out) > 0 {
				stdout.Write([]byte(l.Out))
			}
			if len(l.Err) > 0 {
				stderr.Write([]byte(l.Err))
			}
		}
	}
}
