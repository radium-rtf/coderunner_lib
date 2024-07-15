package container

import (
	"context"
	"errors"
	dc "github.com/docker/docker/api/types/container"
	"github.com/khostya/coderunner/internal/docker"
	"sync"
	"sync/atomic"
)

type Container struct {
	ctx context.Context

	isStarted atomic.Bool

	wg *sync.WaitGroup

	id     string
	docker *docker.Client
}

func NewContainer(ctx context.Context, id string, docker *docker.Client) *Container {
	return &Container{
		id:     id,
		docker: docker,
		ctx:    ctx,
	}
}

func (c *Container) Start(options dc.StartOptions) error {
	if c.isStarted.Swap(true) {
		return ErrContainerHasAlreadyStarted
	}
	return c.docker.StartContainer(c.ctx, c.id, options)
}

func (c *Container) ShowStdout() (string, error) {
	return c.docker.ShowStdoutContainer(c.ctx, c.id)
}

func (c *Container) Wait() (int64, error) {
	statusCh, errCh := c.docker.ContainerWait(c.ctx, c.id, dc.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return 0, err
	case code := <-statusCh:
		if code.Error != nil {
			return 0, errors.New(code.Error.Message)
		}
		return code.StatusCode, nil
	}
}

func (c *Container) Close() error {
	return c.docker.ForceDeleteContainer(c.ctx, c.id)
}
