package docker

import (
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/radium-rtf/coderunner_lib/internal/info"
	"io"
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
	return logs.Stdout, nil
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
