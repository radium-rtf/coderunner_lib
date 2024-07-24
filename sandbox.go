package coderunner

import (
	dc "github.com/docker/docker/api/types/container"
	"github.com/radium-rtf/coderunner_lib/internal/container"
)

type Sandbox struct {
	container *container.Container
}

func newSandbox(container *container.Container) *Sandbox {
	return &Sandbox{container: container}
}

func (s Sandbox) Start() error {
	return s.container.Start(dc.StartOptions{})
}

func (s Sandbox) Wait() (int64, error) {
	return s.container.Wait()
}

func (s Sandbox) ShowStdout() (string, error) {
	return s.container.ShowStdout()
}

func (s Sandbox) ShowStd() (Logs, error) {
	logs, err := s.container.ShowStd()
	if err != nil {
		return Logs{}, err
	}

	return newLogs(logs), nil
}

func (s Sandbox) Close() error {
	return s.container.Close()
}
