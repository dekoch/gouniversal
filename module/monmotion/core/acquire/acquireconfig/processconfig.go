package acquireconfig

type ProcessConfig struct {
	Resolution
	Crop bool
}

func (hc *ProcessConfig) Lock() {

	mut.Lock()
}

func (hc *ProcessConfig) Unlock() {

	mut.Unlock()
}

func (hc *ProcessConfig) LoadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	hc.Width = 0
	hc.Height = 0
}

func (hc *ProcessConfig) SetResolution(res Resolution) {

	mut.Lock()
	defer mut.Unlock()

	hc.Resolution = res
}

func (hc *ProcessConfig) GetResolution() Resolution {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Resolution
}

func (hc *ProcessConfig) SetCrop(enable bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.Crop = enable
}

func (hc *ProcessConfig) GetCrop() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Crop
}
