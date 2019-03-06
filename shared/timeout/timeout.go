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

func (to *TimeOut) Start(millis int) {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.setTimeOut(millis)
	to.enable(true)
	to.reset()
}

func (to *TimeOut) SetTimeOut(millis int) {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.setTimeOut(millis)
}

func (to *TimeOut) setTimeOut(millis int) {

	to.max = time.Duration(millis) * time.Millisecond
}

func (to *TimeOut) Enable(b bool) {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.enable(b)
}

func (to *TimeOut) enable(b bool) {

	to.enabled = b
}

func (to *TimeOut) Reset() {

	to.mut.Lock()
	defer to.mut.Unlock()

	to.reset()
}

func (to *TimeOut) reset() {

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

func (to *TimeOut) ElapsedChan() <-chan bool {

	out := make(chan bool)

	go func() {

		defer close(out)

		for {
			to.mut.RLock()

			if to.enabled == false {
				to.mut.RUnlock()
				return
			}

			t := time.Since(to.t)
			if t > to.max {
				to.mut.RUnlock()
				return
			}

			to.mut.RUnlock()

			time.Sleep(10 * time.Millisecond)
		}
	}()

	return out
}

func (to *TimeOut) ElapsedMillis() float64 {

	to.mut.RLock()
	defer to.mut.RUnlock()

	if to.enabled == false {
		return 0.0
	}

	t := time.Since(to.t)

	return t.Seconds() * 1000.0
}
