package limit

type Opt func(*Limits)

func WithTimeoutInSec(sec int) Opt {
	return func(limits *Limits) {
		limits.TimeoutInSec = &sec
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
