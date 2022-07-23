package runtime

import (
	"io"
	"os"
)

const (
	DetachKeys = "ctrl-d"
)

type Streams struct {
	In         io.Reader
	Out        io.Writer
	Err        io.Writer
	DetachKeys string
}

var StreamsStdio = NewStreams(os.Stdin, os.Stdout, os.Stderr)

func NewStreams(stdin io.Reader, stdout, stderr io.Writer) *Streams {
	return &Streams{
		In:         stdin,
		Out:        stdout,
		Err:        stderr,
		DetachKeys: DetachKeys,
	}
}
