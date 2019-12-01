package acquire

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/webcam"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"

	"github.com/disintegration/imaging"
)

type Acquire struct {
	useWebcam          bool
	webserverChecked   time.Time
	webcam             webcam.Webcam
	Source             string
	Width, Height, Fps uint32
}

func (ac *Acquire) LoadDefaults() {

	ac.Source = "/dev/video0"
	ac.Width = 0
	ac.Height = 0
	ac.Fps = 30
}

func (ac *Acquire) Start() error {

	if strings.HasPrefix(ac.Source, "http") {

		ac.useWebcam = false
	} else {
		ac.useWebcam = true

		return ac.webcam.Start(ac.Source, ac.Width, ac.Height, ac.Fps)
	}

	return nil
}

func (ac *Acquire) GetImage() (typemd.MoImage, error) {

	if ac.useWebcam {
		return ac.webcam.GetImage()
	}

	return ac.fromWebsever()
}

func (ac *Acquire) fromWebsever() (typemd.MoImage, error) {

	pause := time.Duration(1000/ac.Fps)*time.Millisecond - time.Since(ac.webserverChecked)

	time.Sleep(pause)

	var (
		err error
		ret typemd.MoImage
		b   []byte
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				b, err = download(ac.Source)

			case 1:
				ret.Captured = time.Now()
				ret.Img, err = imaging.Decode(bytes.NewReader(b), imaging.AutoOrientation(true))

			case 2:
				ret.Img = imaging.Resize(ret.Img, int(ac.Width), int(ac.Height), imaging.Box)

				ac.webserverChecked = time.Now()
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func download(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
