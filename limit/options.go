package limit

import "time"

type Opt func(*Limits)

func WithTimeout(timeout time.Duration) Opt {
	return func(limits *Limits) {
		limits.Timeout = timeout
	}
}

func WithCPUCount(count int64) Opt {
	return func(limits *Limits) {
		limits.CPUCount = count
	}
}

func WithMemoryInBytes(bytes int64) Opt {
	return func(limits *Limits) {
		limits.MemoryInBytes = bytes
	}
}
