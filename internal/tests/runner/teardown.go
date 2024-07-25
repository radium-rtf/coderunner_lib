package runner

import "github.com/docker/docker/client"

func (r Runner) TearDown() {
	_, _ = client.NewClientWithOpts(client.WithHost(r.dockerHost))
}
