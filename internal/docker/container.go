package docker

import (
	"bytes"
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
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
	reader, err := c.logsContainer(ctx, id, dc.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = stdcopy.StdCopy(&buf, &bytes.Buffer{}, reader)
	return buf.String(), err
}

func (c Client) logsContainer(ctx context.Context, id string, options dc.LogsOptions) (io.ReadCloser, error) {
	reader, err := c.docker.ContainerLogs(ctx, id, options)
	return reader, err
}

func (c Client) ContainerWait(ctx context.Context, id string, running dc.WaitCondition) (<-chan dc.WaitResponse, <-chan error) {
	return c.docker.ContainerWait(ctx, id, running)
}
