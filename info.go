package coderunner

import "github.com/radium-rtf/coderunner_lib/internal/info"

type Logs struct {
	StdOut, StdErr string
}

func newLogs(logs info.Logs) Logs {
	return Logs{
		StdOut: logs.Stdout,
		StdErr: logs.Stderr,
	}
}
