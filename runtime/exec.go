package runtime

import (
	"io"
	"os"
)

type Exec struct {
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Terminal bool
}

func ExecStdio() *Exec {
	return &Exec{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func ExecStreams(stdin io.Reader, stdout, stderr io.Writer) *Exec {
	return &Exec{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}
