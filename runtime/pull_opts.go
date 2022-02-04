package runtime

import (
	"io"
	"os"
)

type Pull struct {
	Stdout io.Writer
}

type PullOpt func(po *Pull) error

func NewPull(opts ...PullOpt) (*Pull, error) {
	p := &Pull{
		Stdout: os.Stdout,
	}
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func WithStream(s io.Writer) PullOpt {
	return func(po *Pull) error {
		po.Stdout = s
		return nil
	}
}
