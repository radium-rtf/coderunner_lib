package docker

import (
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/radium-rtf/coderunner_lib/info"
	"github.com/radium-rtf/coderunner_lib/internal/status"
	"io"
	"time"
)

func (c Client) CreateContainer(ctx context.Context, config *dc.Config, hostConfig *dc.HostConfig,
	networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (dc.CreateResponse, error) {
	resp, err := c.docker.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
	return resp, err
}

func (c Client) ForceDeleteContainer(ctx context.Context, id string) error {
	return c.deleteContainer(ctx, id, dc.RemoveOptions{Force: true})
}

func (c Client) deleteContainer(ctx context.Context, id string, options dc.RemoveOptions) error {
	err := c.docker.ContainerRemove(ctx, id, options)
	return err
}

func (c Client) StartContainer(ctx context.Context, id string, options dc.StartOptions) error {
	err := c.docker.ContainerStart(ctx, id, options)
	return err
}

func (c Client) ShowStdoutContainer(ctx context.Context, id string) (string, error) {
	logs, err := c.ShowStdContainer(ctx, id, dc.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	return logs.StdOut, nil
}

func (c Client) ShowStdContainer(ctx context.Context, id string, options dc.LogsOptions) (info.Logs, error) {
	reader, err := c.logsContainer(ctx, id, options)
	if err != nil {
		return info.Logs{}, err
	}
	defer reader.Close()

	return info.NewLogs(reader)
}

func (c Client) logsContainer(ctx context.Context, id string, options dc.LogsOptions) (io.ReadCloser, error) {
	return c.docker.ContainerLogs(ctx, id, options)
}

func (c Client) ContainerWait(ctx context.Context, id string, running dc.WaitCondition) (<-chan dc.WaitResponse, <-chan error) {
	return c.docker.ContainerWait(ctx, id, running)
}

func (c Client) ContainerFinishedAt(ctx context.Context, id string) (time.Time, error) {
	inspect, err := c.docker.ContainerInspect(ctx, id)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, inspect.State.FinishedAt)
}

func (c Client) ContainerOOMKilled(ctx context.Context, id string) (bool, error) {
	status, err := c.ContainerStatus(ctx, id)
	if err != nil {
		return false, err
	}
	return status.OOMKilled, nil
}

func (c Client) ContainerStatus(ctx context.Context, id string) (status.Status, error) {
	inspect, err := c.docker.ContainerInspect(ctx, id)
	if err != nil {
		return status.Status{}, err
	}

	status := status.Status{
		OOMKilled: inspect.State.OOMKilled,
		Code:      int64(inspect.State.ExitCode),
		Status:    status.T(inspect.State.Status),
	}

	if inspect.State.FinishedAt != "" {
		status.Time.FinishedAt, err = time.Parse(time.RFC3339, inspect.State.FinishedAt)
	}
	if err != nil {
		return status, err
	}

	if inspect.State.StartedAt != "" {
		status.Time.StartedAt, err = time.Parse(time.RFC3339, inspect.State.StartedAt)
	}
	if err != nil {
		return status, err
	}

	return status, nil
}
