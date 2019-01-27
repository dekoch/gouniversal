package sint

import "sync"

type Sint struct {
	mut sync.RWMutex
	val int
}

func (si *Sint) Set(i int) {

	si.mut.Lock()
	defer si.mut.Unlock()

	si.val = i
}

func (si *Sint) Get() int {

	si.mut.RLock()
	defer si.mut.RUnlock()

	return si.val
}
