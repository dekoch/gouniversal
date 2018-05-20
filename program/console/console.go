package console

import "sync"

type Console struct {
	Mut       sync.Mutex
	UserInput string
}

func (c *Console) Input(s string) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	c.UserInput = s
}

func (c *Console) Get() string {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	return c.UserInput
}
