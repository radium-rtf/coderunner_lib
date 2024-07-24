package info

import (
	"bytes"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
)

type Logs struct {
	Stdout string
	Stderr string
}

func NewLogs(reader io.ReadCloser) (Logs, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, err := stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		return Logs{}, err
	}
	return Logs{
		Stderr: stderr.String(),
		Stdout: stdout.String(),
	}, nil
}
