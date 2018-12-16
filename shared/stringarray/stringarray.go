package stringarray

import "sync"

type StringArray struct {
	mut sync.RWMutex
	str []string
}

func (s *StringArray) Add(str string) {

	s.mut.Lock()
	defer s.mut.Unlock()

	s.str = append(s.str, str)
}

func (s *StringArray) AddList(list []string) {

	for _, l := range list {

		s.Add(l)
	}
}

func (s *StringArray) List() []string {

	s.mut.RLock()
	defer s.mut.RUnlock()

	return s.str
}

func (s *StringArray) Count() int {

	s.mut.RLock()
	defer s.mut.RUnlock()

	return len(s.str)
}

func (s *StringArray) Remove(str string) {

	s.mut.Lock()
	defer s.mut.Unlock()

	var l []string

	for i := 0; i < len(s.str); i++ {

		if str != s.str[i] {

			l = append(l, s.str[i])
		}
	}

	s.str = l
}

func (s *StringArray) RemoveAll() {

	s.mut.Lock()
	defer s.mut.Unlock()

	var l []string
	s.str = l
}
