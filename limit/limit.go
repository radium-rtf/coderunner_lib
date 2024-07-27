package limit

import "time"

type Limits struct {
	Timeout       time.Duration
	MemoryInBytes int64
	CPUCount      int64
}

func NewLimits(opts ...Opt) *Limits {
	defaultCfg := &Limits{
		Timeout:       time.Second * 30,
		MemoryInBytes: 1024 * 1024 * 20,
		CPUCount:      1,
	}

	for _, opt := range opts {
		opt(defaultCfg)
	}

	return defaultCfg
}
