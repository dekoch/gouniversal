package webcam

import (
	"bytes"
	"errors"
	"image/jpeg"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/mjpegavi1"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
)

type Webcam struct {
	source             string
	width, height, fps int
	dev                *v4l.Device
}

func (we *Webcam) Start(source string, width, height, fps uint32) error {

	we.source = source
	we.width = int(width)
	we.height = int(height)
	we.fps = int(fps)

	var (
		err error
		cfg v4l.DeviceConfig
	)

	for i := 0; i <= 7; i++ {

		switch i {
		case 0:
			devs := v4l.FindDevices()
			if len(devs) < 1 {
				return errors.New("no device found")
			}

			found := false

			for _, dev := range devs {
				if dev.Path == we.source {
					found = true
				}
			}

			if found == false {
				return errors.New("device " + we.source + " not found")
			}

		case 1:
			we.dev, err = v4l.Open(we.source)

		case 2:
			cfg, err = we.dev.GetConfig()

		case 3:
			cfg.Format = mjpeg.FourCC

			if we.width > 0 {
				cfg.Width = we.width
			}

			if we.height > 0 {
				cfg.Height = we.height
			}

			if we.fps > 0 {
				cfg.FPS.N = uint32(we.fps)
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
			_, err = we.GetImage()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (we *Webcam) GetImage() (typemd.MoImage, error) {

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
