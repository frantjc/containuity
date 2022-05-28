package grpcio

import (
	"fmt"
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
			switch l.Stream {
			case 0:
				stdout.Write(l.Data)
			case 1:
				stderr.Write(l.Data)
			default:
				return fmt.Errorf("unknown stream '%d', must be '0' or '1' for stdout or stderr, respectively", l.Stream)
			}
		}
	}
}
