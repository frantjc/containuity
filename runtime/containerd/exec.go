package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/cio"
	"github.com/frantjc/sequence/runtime"
)

func (c *containerdContainer) Exec(ctx context.Context, opts ...runtime.ExecOpt) error {
	e, err := runtime.NewExec(opts...)
	if err != nil {
		return err
	}

	iopts := []cio.Opt{cio.WithStreams(e.Stdin, e.Stdout, e.Stderr)}
	if e.Terminal {
		iopts = append(iopts, cio.WithTerminal)
	}
	task, err := c.container.NewTask(ctx, cio.NewCreator(iopts...))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return err
	}

	if err := task.Start(ctx); err != nil {
		return err
	}

	exitStatus := <-exitStatusC
	if exitCode, _, err := exitStatus.Result(); err != nil {
		return err
	} else if exitCode != 0 {
		return fmt.Errorf("job exited with code %d", exitCode)
	}

	return nil
}
