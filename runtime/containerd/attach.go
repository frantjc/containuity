package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/cio"
	"github.com/frantjc/sequence/runtime"
)

func (c *containerdContainer) Attach(ctx context.Context, opts ...runtime.ExecOpt) error {
	e, err := runtime.NewExec(opts...)
	if err != nil {
		return err
	}

	iopts := []cio.Opt{cio.WithStreams(e.Stdin, e.Stdout, e.Stderr)}
	if e.Terminal {
		iopts = append(iopts, cio.WithTerminal)
	}
	task, err := c.container.Task(ctx, cio.NewAttach(iopts...))
	if err != nil {
		return err
	}

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return err
	}

	select {
	case exitStatus := <-exitStatusC:
		if exitCode, _, err := exitStatus.Result(); err != nil {
			return err
		} else if exitCode != 0 {
			return fmt.Errorf("job exited with code %d", exitCode)
		}
	case <-ctx.Done():
	}

	return nil
}
