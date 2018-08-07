package stringArray

import "sync"

type StringArray struct {
	mut sync.RWMutex
	str []string
}

func (s *StringArray) Add(str string) {

	s.mut.Lock()
	defer s.mut.Unlock()

	newString := make([]string, 1)
	newString[0] = str

	s.str = append(s.str, newString...)
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
	n := make([]string, 1)

	for i := 0; i < len(s.str); i++ {

		if str != s.str[i] {

			n[0] = s.str[i]

			l = append(l, n...)
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
