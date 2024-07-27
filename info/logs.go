package info

import (
	"bytes"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
)

type Logs struct {
	StdOut, StdErr string
}

func NewLogs(reader io.ReadCloser) (Logs, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, err := stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		return Logs{}, err
	}
	return Logs{
		StdErr: stderr.String(),
		StdOut: stdout.String(),
	}, nil
}
