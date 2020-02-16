package webcam

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"

	"github.com/jinzhu/copier"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
)

type Webcam struct {
	config      acquireconfig.DeviceConfig
	dev         *v4l.Device
	devConfig   v4l.DeviceConfig
	listConfigs []v4l.DeviceConfig
	active      bool
}

var mut sync.RWMutex

func (we *Webcam) Test(conf acquireconfig.DeviceConfig) error {

	err := we.Start(conf)
	if err != nil {
		return err
	}

	return we.Stop()
}

func (we *Webcam) Start(conf acquireconfig.DeviceConfig) error {

	mut.Lock()
	defer mut.Unlock()

	if we.active {
		return nil
	}

	we.config.Lock()
	copier.Copy(&we.config, &conf)
	we.config.Unlock()

	var (
		err error
		cfg v4l.DeviceConfig
	)

	for i := 0; i <= 7; i++ {

		switch i {
		case 0:
			if IsDeviceAvailable(conf.Source) == false {
				return errors.New("device " + conf.Source + " not found")
			}

		case 1:
			we.dev, err = v4l.Open(conf.Source)

		case 2:
			we.active = true

			we.listConfigs, err = we.dev.ListConfigs()

			if len(we.listConfigs) > 0 {
				cfg = we.listConfigs[0]
			}

		case 3:
			cfg.Format = mjpeg.FourCC

			if we.config.Width > 0 {
				cfg.Width = we.config.Width
			}

			if we.config.Height > 0 {
				cfg.Height = we.config.Height
			}

			if we.config.FPS > 0 {
				cfg.FPS.N = uint32(we.config.FPS)
				cfg.FPS.D = 1
			}

			err = we.dev.SetConfig(cfg)

		case 4:
			err = we.dev.TurnOn()

		case 5:
			we.devConfig, err = we.dev.GetConfig()

		case 6:
			if we.devConfig.Format != mjpeg.FourCC {
				return errors.New("format not supported")
			}

		case 7:
			var img mdimg.MDImage
			err = we.getImage(&img)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (we *Webcam) Stop() error {

	mut.Lock()
	defer mut.Unlock()

	if we.active == false {
		return nil
	}

	if IsDeviceAvailable(we.config.Source) == false {
		return errors.New("device " + we.config.Source + " not found")
	}

	we.dev.Close()

	we.active = false

	return nil
}

func (we *Webcam) GetImage(image *mdimg.MDImage) error {

	mut.RLock()
	defer mut.RUnlock()

	return we.getImage(image)
}

func (we *Webcam) getImage(image *mdimg.MDImage) error {

	var err error

	func() {

		var (
			buf *v4l.Buffer
			b   []byte
		)

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				buf, err = we.dev.Capture()

			case 1:
				b = make([]byte, buf.Size())
				buf.ReadAt(b, 0)

			case 2:
				image.Jpeg = b
				image.Captured = time.Now()
				image.Width = we.devConfig.Width
				image.Height = we.devConfig.Height
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (we *Webcam) ListConfigs() ([]acquireconfig.DeviceConfig, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret []acquireconfig.DeviceConfig
	)

	for i := 0; i <= 1; i++ {

		switch i {
		case 0:
			if we.active == false {
				err = we.devListConfigs()
			}

		case 1:
			if len(we.listConfigs) == 0 {
				return ret, nil
			}

			for i := range we.listConfigs {

				if we.listConfigs[i].Format != mjpeg.FourCC {
					continue
				}

				var n acquireconfig.DeviceConfig
				n.Width = we.listConfigs[i].Width
				n.Height = we.listConfigs[i].Height
				n.FPS = int(we.listConfigs[i].FPS.N)

				ret = append(ret, n)
			}
		}

		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func (we *Webcam) devListConfigs() error {

	if we.active {
		return errors.New("device is active")
	}

	var err error

	for i := 0; i <= 3; i++ {

		switch i {
		case 0:
			if IsDeviceAvailable(we.config.Source) == false {
				return errors.New("device " + we.config.Source + " not found")
			}

		case 1:
			we.dev, err = v4l.Open(we.config.Source)

		case 2:
			we.listConfigs, err = we.dev.ListConfigs()

		case 3:
			we.dev.Close()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func FindDevices() []string {

	var ret []string

	for _, dev := range v4l.FindDevices() {
		ret = append(ret, dev.Path)
	}

	return ret
}

func IsDeviceAvailable(path string) bool {

	devs := v4l.FindDevices()
	if len(devs) == 0 {
		return false
	}

	for _, dev := range devs {
		if dev.Path == path {
			return true
		}
	}

	return false
}
