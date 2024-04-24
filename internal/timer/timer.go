package timer

import "time"

type TimerState int

const (
	Stopped TimerState = iota
)

type Timer struct {
	start time.Time
}

func NewTimer() Timer {
	return Timer{
		start: time.Now(),
	}
}

func (t *Timer) Begin() {
	t.start = time.Now()
}

// / Returns the curernt time in ms
func (t Timer) CurrentTime() int64 {
	return time.Since(t.start).Milliseconds()
}
