package container

import (
	"context"
	"errors"
	dc "github.com/docker/docker/api/types/container"
	"github.com/radium-rtf/coderunner_lib/info"
	"github.com/radium-rtf/coderunner_lib/internal/docker"
	"github.com/radium-rtf/coderunner_lib/internal/status"
	"sync/atomic"
	"time"
)

type result struct {
	code int64
	err  error
}

type Container struct {
	ctx context.Context

	isStarted atomic.Bool

	id     string
	docker *docker.Client

	timeLimit time.Duration

	wait   chan struct{}
	result result
}

func NewContainer(ctx context.Context, id string, docker *docker.Client, timeout time.Duration) *Container {
	return &Container{
		ctx:       ctx,
		id:        id,
		docker:    docker,
		timeLimit: timeout,
		wait:      make(chan struct{}),
	}
}

func (c *Container) Start(options dc.StartOptions) error {
	if c.isStarted.Swap(true) {
		return ErrContainerHasAlreadyStarted
	}

	err := c.docker.StartContainer(c.ctx, c.id, options)
	if err != nil {
		return err
	}

	timeout, _ := context.WithTimeout(c.ctx, c.timeLimit)
	go c.startWait(timeout)
	return nil
}

func (c *Container) startWait(timeout context.Context) {
	defer close(c.wait)

	var (
		code int64
		err  error
	)

	select {
	case <-timeout.Done():
		code, err = 0, ErrContainerTimeout
	default:
		break
	}

	statusCh, errCh := c.docker.ContainerWait(c.ctx, c.id, dc.WaitConditionNotRunning)
	select {
	case <-timeout.Done():
		code, err = 0, ErrContainerTimeout
	case e := <-errCh:
		c.result = result{code: code, err: e}
		return
	case status := <-statusCh:
		if status.Error != nil {
			c.result = result{code: code, err: errors.New(status.Error.Message)}
			return
		}
		c.result = result{code: status.StatusCode}
		return
	}

	select {
	case e := <-errCh:
		code, err = 0, e
	case status := <-statusCh:
		if status.Error != nil {
			code, err = 0, errors.New(status.Error.Message)
		}
		code, err = status.StatusCode, nil
	default:
		break
	}
	c.result = result{code: code, err: err}
}

func (c *Container) ShowStdout() (string, error) {
	return c.docker.ShowStdoutContainer(c.ctx, c.id)
}

func (c *Container) ShowStd() (info.Logs, error) {
	return c.docker.ShowStdContainer(c.ctx, c.id, dc.LogsOptions{ShowStderr: true, ShowStdout: true})
}

func (c *Container) FinishedAt() (time.Time, error) {
	return c.docker.ContainerFinishedAt(c.ctx, c.id)
}

func (c *Container) Status() (status.Status, error) {
	return c.docker.ContainerStatus(c.ctx, c.id)
}

func (c *Container) Wait() (int64, error) {
	if !c.isStarted.Load() {
		return 0, ErrContainerHasNotBeenStarted
	}

	<-c.wait
	return c.result.code, c.result.err
}

func (c *Container) Close() error {
	return c.docker.ForceDeleteContainer(c.ctx, c.id)
}
