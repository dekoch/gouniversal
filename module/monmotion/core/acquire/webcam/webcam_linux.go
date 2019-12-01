// +build linux

package webcam

/*
#include <linux/videodev2.h>
*/
import "C"
import (
	"bytes"
	"errors"
	"image"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
)

var (
	dhtMarker = []byte{255, 196}
	dht       = []byte{1, 162, 0, 0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 16, 0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125, 1, 2, 3, 0, 4, 17, 5, 18, 33, 49, 65, 6, 19, 81, 97, 7, 34, 113, 20, 50, 129, 145, 161, 8, 35, 66, 177, 193, 21, 82, 209, 240, 36, 51, 98, 114, 130, 9, 10, 22, 23, 24, 25, 26, 37, 38, 39, 40, 41, 42, 52, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 17, 0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119, 0, 1, 2, 3, 17, 4, 5, 33, 49, 6, 18, 65, 81, 7, 97, 113, 19, 34, 50, 129, 8, 20, 66, 145, 161, 177, 193, 9, 35, 51, 82, 240, 21, 98, 114, 209, 10, 22, 36, 52, 225, 37, 241, 23, 24, 25, 26, 38, 39, 40, 41, 42, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 130, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 226, 227, 228, 229, 230, 231, 232, 233, 234, 242, 243, 244, 245, 246, 247, 248, 249, 250}
	sosMarker = []byte{255, 218}
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

	for i := 0; i <= 6; i++ {

		switch i {
		case 0:
			devs := v4l.FindDevices()
			if len(devs) < 1 {
				return errors.New("no device found")
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

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				buf, err = we.dev.Capture()

			case 1:
				b = make([]byte, buf.Size())
				buf.ReadAt(b, 0)

			case 2:
				ret.Captured = time.Now()

				var s string
				ret.Img, s, err = image.Decode(bytes.NewReader(b))

				if s != "jpeg" {
					err = errors.New("unsupported format " + s)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}
