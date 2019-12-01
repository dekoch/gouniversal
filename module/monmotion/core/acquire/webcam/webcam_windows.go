// +build windows

package webcam

import (
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
)

type Webcam struct {
	source             string
	width, height, fps uint32
}

func (we *Webcam) Start(source string, width, height, fps uint32) error {

	we.source = source
	we.width = width
	we.height = height
	we.fps = fps

	return nil
}

func (we *Webcam) GetImage() (typemd.MoImage, error) {

	var (
		err error
		ret typemd.MoImage
	)

	return ret, err
}
