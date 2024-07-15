package docker

import (
	"github.com/docker/docker/client"
)

type Client struct {
	docker  *client.Client
	workDir string
	uid     int
}

func NewClient(wordDir string, uid int, opts ...client.Opt) (*Client, error) {
	docker, err := client.NewClientWithOpts(opts...)

	if err != nil {
		return nil, err
	}
	return &Client{docker: docker, workDir: wordDir, uid: uid}, nil
}

func (c Client) Close() error {
	return c.docker.Close()
}
