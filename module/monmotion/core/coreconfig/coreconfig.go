package coreconfig

import (
	"sync"
	"time"
)

type CoreConfig struct {
	Name        string
	Enabled     bool
	Record      bool
	FileRoot    string
	PreRecoding int // millisecond
	Overrun     int // millisecond
	Setup       int // millisecond
	mut         sync.RWMutex
}

func (hc *CoreConfig) LoadCoreDefaults() {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Name = "motion0"
	hc.Enabled = true
	hc.Record = false
	hc.FileRoot = "data/monmotion/"
	hc.PreRecoding = 10000
	hc.Overrun = 10000
	hc.Setup = 3000
}

func (hc *CoreConfig) SetName(name string) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Name = name
}

func (hc *CoreConfig) GetName() string {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.Name
}

func (hc *CoreConfig) SetRecord(record bool) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Record = record
}

func (hc *CoreConfig) GetRecord() bool {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.Record
}

func (hc *CoreConfig) SetFileRoot(path string) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.FileRoot = path
}

func (hc *CoreConfig) GetFileRoot() string {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.FileRoot
}

func (hc *CoreConfig) SetSetup(millis int) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Setup = millis
}

func (hc *CoreConfig) GetSetup() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return time.Duration(hc.Setup) * time.Millisecond
}

func (hc *CoreConfig) SetPreRecoding(millis int) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.PreRecoding = millis
}

func (hc *CoreConfig) GetPreRecoding() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return time.Duration(hc.PreRecoding) * time.Millisecond
}

func (hc *CoreConfig) SetOverrun(millis int) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Overrun = millis
}

func (hc *CoreConfig) GetOverrun() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return time.Duration(hc.Overrun) * time.Millisecond
}

func (hc *CoreConfig) GetRecodingDuration() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return time.Duration(hc.PreRecoding+hc.Overrun) * time.Millisecond
}
