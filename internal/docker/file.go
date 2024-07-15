package docker

import (
	"context"
	dc "github.com/docker/docker/api/types/container"
	"github.com/khostya/coderunner/file"
	"github.com/khostya/coderunner/internal/tar"
)

func (c Client) WriteFiles(ctx context.Context, id string, files []file.File) error {
	reader, err := tar.NewTarArchive(files, c.uid)
	if err != nil {
		return err
	}

	err = c.docker.CopyToContainer(ctx, id, c.workDir, reader, dc.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return err
}
