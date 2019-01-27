package timeout

import (
	"sync"
	"time"
)

type TimeOut struct {
	mut     sync.RWMutex
	t       time.Time
	max     time.Duration
	enabled bool
}

func (to *TimeOut) SetTimeOut(millis int) {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.max = time.Duration(millis) * time.Millisecond
}

func (to *TimeOut) Enable(b bool) {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.enabled = b
}

func (to *TimeOut) Reset() {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.t = time.Now()
}

func (to *TimeOut) Elapsed() bool {

	to.mut.RLock()
	defer to.mut.RUnlock()

	if to.enabled == false {
		return false
	}

	t := time.Since(to.t)
	if t < to.max {
		return false
	}

	return true
}
