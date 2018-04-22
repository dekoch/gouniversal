package programTypes

import (
	"sync"
)

type Console struct {
	Mut   sync.Mutex
	Input string
}
