package coreconfig

import (
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"

	"github.com/google/uuid"
)

type CoreConfig struct {
	UUID        string
	Name        string
	Enabled     bool
	Record      bool
	FileRoot    string
	PreRecoding int // millisecond
	Overrun     int // millisecond
	Setup       int // millisecond
	Acquire     acquireconfig.AcquireConfig
	Trigger     triggerconfig.TriggerConfig
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
	hc.PreRecoding = 3000
	hc.Overrun = 3000
	hc.Setup = 3000
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

func (hc *CoreConfig) SetSetup(millis int) {

	mut.Lock()
	defer mut.Unlock()

	hc.Setup = millis
}

func (hc *CoreConfig) GetSetup() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.Setup) * time.Millisecond
}

func (hc *CoreConfig) SetPreRecoding(millis int) {

	mut.Lock()
	defer mut.Unlock()

	hc.PreRecoding = millis
}

func (hc *CoreConfig) GetPreRecoding() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.PreRecoding) * time.Millisecond
}

func (hc *CoreConfig) SetOverrun(millis int) {

	mut.Lock()
	defer mut.Unlock()

	hc.Overrun = millis
}

func (hc *CoreConfig) GetOverrun() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.Overrun) * time.Millisecond
}

func (hc *CoreConfig) GetRecodingDuration() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.PreRecoding+hc.Overrun) * time.Millisecond
}
