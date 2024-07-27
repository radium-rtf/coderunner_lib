package coderunner

import (
	"errors"
	dc "github.com/docker/docker/api/types/container"
	"github.com/radium-rtf/coderunner_lib/info"
	"github.com/radium-rtf/coderunner_lib/internal/container"
	"github.com/radium-rtf/coderunner_lib/internal/status"
	"time"
)

type Sandbox struct {
	container *container.Container
	timeout   time.Duration
}

func newSandbox(container *container.Container, timeout time.Duration) *Sandbox {
	return &Sandbox{container: container, timeout: timeout}
}

func (s Sandbox) Start() error {
	return s.container.Start(dc.StartOptions{})
}

func (s Sandbox) Wait() (int64, error) {
	status, err := s.container.Wait()
	if err == nil {
		return status, nil
	}

	if errors.Is(err, container.ErrContainerTimeout) {
		return 0, ErrTimeout
	}
	return 0, err
}

func (s Sandbox) ShowStdout() (string, error) {
	return s.container.ShowStdout()
}

func (s Sandbox) ShowStd() (info.Logs, error) {
	logs, err := s.container.ShowStd()
	if err != nil {
		return info.Logs{}, err
	}

	return logs, nil
}

func (s Sandbox) Info() (info.Info, error) {
	code, waitErr := s.Wait()
	if waitErr != nil && !errors.Is(waitErr, ErrTimeout) {
		return info.Info{}, waitErr
	}
	if errors.Is(waitErr, ErrTimeout) {
		return info.Info{StatusCode: code, Status: info.StatusTimeout}, nil
	}

	logs, err := s.ShowStd()
	if err != nil {
		return info.Info{}, err
	}

	status, err := s.container.Status()
	if err != nil {
		return info.Info{}, err
	}

	return info.Info{
		StatusCode: code,
		Logs:       logs,
		Time:       status.Time,
		Status:     newStatus(status, s.timeout),
	}, nil
}

func (s Sandbox) FinishedAt() (time.Time, error) {
	return s.container.FinishedAt()
}

func (s Sandbox) Close() error {
	return s.container.Close()
}

func newStatus(status status.Status, timeout time.Duration) info.Status {
	switch {
	case status.Time.FinishedAt.Sub(status.Time.StartedAt) > timeout:
		return info.StatusTimeout
	case status.OOMKilled:
		return info.StatusOOMKilled
	case status.Code == 0:
		return info.StatusOK
	default:
		return info.StatusUnknown
	}
}
