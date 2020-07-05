package webcam

import (
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
)

type Webcam struct {
	config acquireconfig.DeviceConfig
}

func (we *Webcam) Test(conf acquireconfig.DeviceConfig) error {

	err := we.Start(conf)
	if err != nil {
		return err
	}

	return we.Stop()
}

func (we *Webcam) Start(conf acquireconfig.DeviceConfig) error {

	return nil
}

func (we *Webcam) Stop() error {

	return nil
}

func (we *Webcam) GetImage(image *mdimg.MDImage) error {

	var err error

	return err
}

func (we *Webcam) ListConfigs() ([]acquireconfig.DeviceConfig, error) {

	var ret []acquireconfig.DeviceConfig

	return ret, nil
}

func FindDevices() []string {

	var ret []string

	return ret
}

func IsDeviceAvailable(path string) bool {

	return false
}
