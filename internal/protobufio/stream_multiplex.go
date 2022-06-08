package protobufio

import (
	"fmt"
	"io"

	protobufiov1 "github.com/frantjc/sequence/internal/protobufio/v1"
)

type LogMessage interface {
	GetLog() *protobufiov1.Log
}

type LogStreamClient[T LogMessage] interface {
	Msg() T
	Err() error
	Receive() bool
	Close() error
}

func DemultiplexLogStream[T LogMessage](stream LogStreamClient[T], stdout, stderr io.Writer) error {
	for stream.Receive() {
		var (
			msg = stream.Msg()
			log = msg.GetLog()
		)
		if log != nil {
			switch log.GetStream() {
			case 0:
				stdout.Write(log.GetData())
			case 1:
				stderr.Write(log.GetData())
			default:
				return fmt.Errorf("unknown stream '%d', must be '0' or '1' for stdout or stderr, respectively", log.GetStream())
			}
		}
	}

	return nil
}
