package throttle

import (
	"sync"
	"time"
)

// Throttler .
type Throttler interface {
	Do(f func())
}

// New .
func New(duration time.Duration) Throttler {
	return &throttle{
		duration: duration,
	}
}

type throttle struct {
	duration time.Duration
	once     sync.Once
}

// Do .
func (t *throttle) Do(f func()) {
	t.once.Do(func() {
		reset := func() {
			time.Sleep(t.duration)
			t.once = sync.Once{}
		}
		go reset()
		f()
	})
}
