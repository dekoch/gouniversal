package sbool

import "sync"

type Sbool struct {
	mut   sync.RWMutex
	state bool
}

func (s *Sbool) Set() {

	s.mut.Lock()
	defer s.mut.Unlock()

	s.state = true
}

func (s *Sbool) UnSet() {

	s.mut.Lock()
	defer s.mut.Unlock()

	s.state = false
}

func (s *Sbool) SetState(b bool) {

	s.mut.Lock()
	defer s.mut.Unlock()

	s.state = b
}

func (s *Sbool) IsSet() bool {

	s.mut.RLock()
	defer s.mut.RUnlock()

	return s.state
}
