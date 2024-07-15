package coderunner

import (
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/radium-rtf/coderunner_lib/config"
	"github.com/radium-rtf/coderunner_lib/file"
	"github.com/radium-rtf/coderunner_lib/internal/container"
	"github.com/radium-rtf/coderunner_lib/internal/docker"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"
)

type Runner struct {
	docker         *docker.Client
	defaultUser    string
	defaultWorkDir string
}

func NewRunner(cfg *config.Config, opt ...client.Opt) (*Runner, error) {
	docker, err := docker.NewClient(cfg.Workdir, cfg.Uid, opt...)
	if err != nil {
		return nil, err
	}
	return &Runner{docker: docker, defaultUser: cfg.User, defaultWorkDir: cfg.Workdir}, nil
}

func (r *Runner) NewSandbox(ctx context.Context, cmd []string,
	profile profile.Profile, limits *limit.Limits, files []file.File) (*Sandbox, error) {
	config := &dc.Config{
		Image:       profile.Image,
		StopTimeout: limits.TimeoutInSec,
		User:        r.defaultUser,
		WorkingDir:  r.defaultWorkDir,
		Cmd:         cmd,
	}
	hostConfig := &dc.HostConfig{
		Resources: dc.Resources{CPUCount: limits.CPUCount, Memory: limits.MemoryInBytes},
	}

	response, err := r.docker.CreateContainer(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	err = r.docker.WriteFiles(ctx, response.ID, files)
	if err != nil {
		return nil, err
	}

	container := container.NewContainer(ctx, response.ID, r.docker)
	return newSandbox(container), nil
}

func (r *Runner) Close() error {
	return r.docker.Close()
}
