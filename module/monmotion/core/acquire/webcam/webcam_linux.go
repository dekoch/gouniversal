package webcam

import (
	"bytes"
	"errors"
	"image/jpeg"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/mjpegavi1"

	"github.com/jinzhu/copier"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
)

type Webcam struct {
	config      acquireconfig.DeviceConfig
	dev         *v4l.Device
	listconfigs []v4l.DeviceConfig
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

			we.listconfigs, err = we.dev.ListConfigs()

			if len(we.listconfigs) > 0 {
				cfg = we.listconfigs[0]
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
			cfg, err = we.dev.GetConfig()

		case 6:
			if cfg.Format != mjpeg.FourCC {
				return errors.New("format not supported")
			}

		case 7:
			_, err = we.getImage()
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

func (we *Webcam) GetImage() (typemd.MoImage, error) {

	mut.RLock()
	defer mut.RUnlock()

	return we.getImage()
}

func (we *Webcam) getImage() (typemd.MoImage, error) {

	var (
		err error
		ret typemd.MoImage
		buf *v4l.Buffer
		b   []byte
	)

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				buf, err = we.dev.Capture()

			case 1:
				b = make([]byte, buf.Size())
				buf.ReadAt(b, 0)

			case 2:
				b, err = mjpegavi1.Decode(b)

			case 3:
				ret.Captured = time.Now()
				ret.Img, err = jpeg.Decode(bytes.NewReader(b))
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
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
			if len(we.listconfigs) == 0 {
				return ret, nil
			}

			for i := range we.listconfigs {

				if we.listconfigs[i].Format != mjpeg.FourCC {
					continue
				}

				var n acquireconfig.DeviceConfig
				n.Width = we.listconfigs[i].Width
				n.Height = we.listconfigs[i].Height
				n.FPS = int(we.listconfigs[i].FPS.N)

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
			we.listconfigs, err = we.dev.ListConfigs()

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
