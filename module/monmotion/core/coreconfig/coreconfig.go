package coreconfig

import (
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"

	"github.com/google/uuid"
)

type CoreConfig struct {
	UUID           string
	Name           string
	Enabled        bool
	Record         bool
	FileRoot       string
	PreRecoding    int // second
	Overrun        int // second
	Setup          int // second
	CacheBlockSize int
	Acquire        acquireconfig.AcquireConfig
	Trigger        triggerconfig.TriggerConfig
}

var mut sync.RWMutex

func (hc *CoreConfig) Lock() {

	mut.Lock()
}

func (hc *CoreConfig) Unlock() {

	mut.Unlock()
}

func (hc *CoreConfig) LoadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	u := uuid.Must(uuid.NewRandom())
	hc.UUID = u.String()
	hc.Name = u.String()
	hc.Enabled = false
	hc.Record = false
	hc.FileRoot = "data/monmotion/"
	hc.PreRecoding = 3
	hc.Overrun = 3
	hc.Setup = 3
	hc.CacheBlockSize = 60
	hc.Acquire.LoadDefaults()
	hc.Trigger.LoadDefaults()
}

func (hc *CoreConfig) SetUUID(uid string) {

	mut.Lock()
	defer mut.Unlock()

	hc.UUID = uid
}

func (hc *CoreConfig) GetUUID() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.UUID
}

func (hc *CoreConfig) SetName(name string) {

	mut.Lock()
	defer mut.Unlock()

	hc.Name = name
}

func (hc *CoreConfig) GetName() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Name
}

func (hc *CoreConfig) SetEnabled(enabled bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.Enabled = enabled
}

func (hc *CoreConfig) GetEnabled() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Enabled
}

func (hc *CoreConfig) SetRecord(record bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.Record = record
}

func (hc *CoreConfig) GetRecord() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Record
}

func (hc *CoreConfig) SetFileRoot(path string) {

	mut.Lock()
	defer mut.Unlock()

	hc.FileRoot = path
}

func (hc *CoreConfig) GetFileRoot() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FileRoot
}

func (hc *CoreConfig) SetSetup(sec int) {

	mut.Lock()
	defer mut.Unlock()

	hc.Setup = sec
}

func (hc *CoreConfig) GetSetup() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.Setup) * time.Second
}

func (hc *CoreConfig) SetPreRecoding(sec int) {

	mut.Lock()
	defer mut.Unlock()

	hc.PreRecoding = sec
}

func (hc *CoreConfig) GetPreRecoding() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.PreRecoding) * time.Second
}

func (hc *CoreConfig) SetOverrun(sec int) {

	mut.Lock()
	defer mut.Unlock()

	hc.Overrun = sec
}

func (hc *CoreConfig) GetOverrun() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.Overrun) * time.Second
}

func (hc *CoreConfig) GetRecodingDuration() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.PreRecoding+hc.Overrun) * time.Second
}

func (hc *CoreConfig) SetCacheBlockSize(size int) {

	mut.Lock()
	defer mut.Unlock()

	hc.CacheBlockSize = size
}

func (hc *CoreConfig) GetCacheBlockSize() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.CacheBlockSize
}
