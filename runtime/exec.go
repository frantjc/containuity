package runtime

import (
	"io"
	"os"
)

type Streams struct {
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Terminal bool
}

var StreamsStdio = &Streams{
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
	Stderr: os.Stderr,
}

func NewStreams(stdin io.Reader, stdout, stderr io.Writer) *Streams {
	return &Streams{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}
