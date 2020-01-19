package acquireconfig

import "github.com/dekoch/gouniversal/module/monmotion/typemd"

type DeviceConfig struct {
	Source string
	typemd.Resolution
	FPS int
}

func (hc *DeviceConfig) Lock() {

	mut.Lock()
}

func (hc *DeviceConfig) Unlock() {

	mut.Unlock()
}

func (hc *DeviceConfig) LoadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	hc.Width = 0
	hc.Height = 0
	hc.FPS = 0
}

func (hc *DeviceConfig) SetSource(src string) {

	mut.Lock()
	defer mut.Unlock()

	hc.Source = src
}

func (hc *DeviceConfig) GetSource() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Source
}

func (hc *DeviceConfig) SetResolution(res typemd.Resolution) {

	mut.Lock()
	defer mut.Unlock()

	hc.Resolution = res
}

func (hc *DeviceConfig) GetResolution() typemd.Resolution {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Resolution
}

func (hc *DeviceConfig) SetFPS(fps int) {

	mut.Lock()
	defer mut.Unlock()

	hc.FPS = fps
}

func (hc *DeviceConfig) GetFPS() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FPS
}
