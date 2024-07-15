package coderunner

import (
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/khostya/coderunner/config"
	"github.com/khostya/coderunner/file"
	"github.com/khostya/coderunner/internal/container"
	"github.com/khostya/coderunner/internal/docker"
	"github.com/khostya/coderunner/limit"
	"github.com/khostya/coderunner/profile"
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
