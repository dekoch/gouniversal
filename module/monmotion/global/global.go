package global

import (
	"errors"
	"sync"

	"github.com/dekoch/gouniversal/module/monmotion/core"
	"github.com/dekoch/gouniversal/module/monmotion/moduleconfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config moduleconfig.ModuleConfig
	Lang   language.Language
	Cores  []core.Core
	mut    sync.Mutex
)

func AddCores(n int) {

	for i := 0; i < n; i++ {
		var nc core.Core
		Cores = append(Cores, nc)
	}
}

func IsCoreAvailable(uid string) bool {

	mut.Lock()
	defer mut.Unlock()

	for i := range Cores {

		if Cores[i].GetUUID() == uid {
			return true
		}
	}

	return false
}

func GetCore(uid string) (*core.Core, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := range Cores {

		if Cores[i].GetUUID() == uid {
			return &Cores[i], nil
		}
	}

	var n core.Core
	return &n, errors.New("device uuid not found")
}

func GetFreeCore() (*core.Core, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := range Cores {

		if Cores[i].GetUUID() == "" {
			return &Cores[i], nil
		}
	}

	var n core.Core
	return &n, errors.New("no free core found")
}

func FreeCore(uid string) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range Cores {

		if Cores[i].GetUUID() == uid {

			err := Cores[i].Stop()
			if err != nil {
				return err
			}

			err = Cores[i].Exit()
			if err != nil {
				return err
			}

			return Cores[i].Reset()
		}
	}

	return errors.New("device uuid not found")
}
