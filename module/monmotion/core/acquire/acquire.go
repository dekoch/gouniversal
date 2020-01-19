package acquire

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/webcam"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"

	"github.com/jinzhu/copier"
)

type Acquire struct {
	config           acquireconfig.AcquireConfig
	webcam           webcam.Webcam
	useWebcam        bool
	webserverChecked time.Time
}

var mut sync.RWMutex

func (ac *Acquire) ListConfigs() ([]acquireconfig.DeviceConfig, error) {

	return ac.webcam.ListConfigs()
}

func (ac *Acquire) TestWebcam(conf acquireconfig.AcquireConfig) error {

	mut.Lock()
	defer mut.Unlock()

	ac.config.Lock()
	copier.Copy(&ac.config, &conf)
	ac.config.Unlock()

	return ac.webcam.Test(ac.config.Device)
}

func (ac *Acquire) Start(conf acquireconfig.AcquireConfig) error {

	mut.Lock()
	defer mut.Unlock()

	ac.config.Lock()
	copier.Copy(&ac.config, &conf)
	ac.config.Unlock()

	if strings.HasPrefix(conf.Device.Source, "http") {

		ac.useWebcam = false
	} else {
		ac.useWebcam = true

		return ac.webcam.Start(ac.config.Device)
	}

	return nil
}

func (ac *Acquire) Stop() error {

	mut.Lock()
	defer mut.Unlock()

	if ac.useWebcam {
		return ac.webcam.Stop()
	}

	return nil
}

func (ac *Acquire) GetImage() (mdimg.MDImage, error) {

	mut.RLock()
	defer mut.RUnlock()

	var (
		err error
		ret mdimg.MDImage
	)

	if ac.useWebcam {
		ret, err = ac.webcam.GetImage()
	} else {
		ret, err = ac.fromWebsever()
	}

	return ret, err
}

func (ac *Acquire) fromWebsever() (mdimg.MDImage, error) {

	pause := time.Duration(1000/ac.config.Device.FPS)*time.Millisecond - time.Since(ac.webserverChecked)

	time.Sleep(pause)

	var (
		err error
		ret mdimg.MDImage
		b   []byte
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				b, err = download(ac.config.Device.Source)

			case 1:
				ret.Captured = time.Now()
				ret.Jpeg = b

			case 2:
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
