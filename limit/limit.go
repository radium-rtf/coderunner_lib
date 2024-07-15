package limit

type Limits struct {
	TimeoutInSec  *int
	MemoryInBytes int64
	CPUCount      int64
}

func NewLimits(opts ...Opt) *Limits {
	var sec = 600

	defaultCfg := &Limits{
		TimeoutInSec:  &sec,
		MemoryInBytes: 1024 * 1024 * 20,
		CPUCount:      1,
	}

	for _, opt := range opts {
		opt(defaultCfg)
	}

	return defaultCfg
}
