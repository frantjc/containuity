package runtime

import (
	"io"
	"os"
)

type Exec struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type ExecOpt func(e *Exec) error

func WithStdio(e *Exec) error {
	e.Stdin = os.Stdin
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	return nil
}

func WithStreams(stdin io.Reader, stdout, stderr io.Writer) ExecOpt {
	return func(e *Exec) error {
		e.Stdin = stdin
		e.Stdout = stdout
		e.Stderr = stderr
		return nil
	}
}

func NewExec(opts ...ExecOpt) (*Exec, error) {
	e := &Exec{}
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}
