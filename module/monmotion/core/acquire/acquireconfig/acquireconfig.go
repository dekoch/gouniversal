package acquireconfig

import (
	"sync"

	"github.com/jinzhu/copier"
)

type AcquireConfig struct {
	Device  DeviceConfig
	Process ProcessConfig
}

type Resolution struct {
	Width  int
	Height int
}

var mut sync.RWMutex

func (hc *AcquireConfig) Lock() {

	mut.Lock()
}

func (hc *AcquireConfig) Unlock() {

	mut.Unlock()
}

func (hc *AcquireConfig) LoadDefaults() {

	hc.Device.LoadDefaults()
	hc.Process.LoadDefaults()
}

func (hc *AcquireConfig) SetDeviceConfig(conf DeviceConfig) error {

	mut.Lock()
	defer mut.Unlock()

	copier.Copy(&hc.Device, &conf)

	return nil
}

func (hc *AcquireConfig) GetDeviceConfig() DeviceConfig {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Device
}

func (hc *AcquireConfig) SetProcessConfig(conf ProcessConfig) error {

	mut.Lock()
	defer mut.Unlock()

	hc.Process = conf

	return nil
}

func (hc *AcquireConfig) GetProcessConfig() ProcessConfig {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Process
}
