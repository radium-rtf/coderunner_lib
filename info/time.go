package info

import "time"

type Time struct {
	StartedAt  time.Time
	FinishedAt time.Time
}

func (t Time) Diff() time.Duration {
	return t.FinishedAt.Sub(t.StartedAt)
}
