package info

type Status string

const (
	StatusOK        = "ok"
	StatusOOMKilled = "oom killed"
	StatusTimeout   = "timeout"
	StatusUnknown   = "unknown"
)
