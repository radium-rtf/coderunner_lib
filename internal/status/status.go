package status

import "github.com/radium-rtf/coderunner_lib/info"

type (
	T string

	Status struct {
		Status    T
		Code      int64
		Time      info.Time
		OOMKilled bool
	}
)

const (
	Created    T = "created"
	Running    T = "running"
	Paused     T = "paused"
	Restarting T = "restarting"
	Removing   T = "removing"
	Exited     T = "exited"
	Dead       T = "dead"
)
