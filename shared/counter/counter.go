package counter

import "sync"

type Counter struct {
	mut   sync.RWMutex
	count int
}

func (c *Counter) Add() {

	c.mut.Lock()
	defer c.mut.Unlock()

	c.count++
}

func (c *Counter) AddCount(n int) {

	c.mut.Lock()
	defer c.mut.Unlock()

	c.count += n
}

func (c *Counter) SetCount(n int) {

	c.mut.Lock()
	defer c.mut.Unlock()

	if n < 0 {
		c.count = 0
		return
	}

	c.count = n
}

func (c *Counter) GetCount() int {

	c.mut.RLock()
	defer c.mut.RUnlock()

	return c.count
}

func (c *Counter) Remove() {

	c.mut.Lock()
	defer c.mut.Unlock()

	if c.count == 0 {
		return
	}

	c.count--
}

func (c *Counter) RemoveCount(n int) {

	c.mut.Lock()
	defer c.mut.Unlock()

	if n > c.count {
		c.count = 0
		return
	}

	c.count -= n
}

func (c *Counter) Reset() {

	c.mut.Lock()
	defer c.mut.Unlock()

	c.count = 0
}
