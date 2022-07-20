package runtime

import (
	"io"
	"os"
)

type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

var StreamsStdio = &Streams{
	In:  os.Stdin,
	Out: os.Stdout,
	Err: os.Stderr,
}

func NewStreams(stdin io.Reader, stdout, stderr io.Writer) *Streams {
	return &Streams{
		In:  stdin,
		Out: stdout,
		Err: stderr,
	}
}
