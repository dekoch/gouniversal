package core

import (
	"errors"
	"image"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire"
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/webcam"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/imgcache"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/uirequest"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"

	"github.com/jinzhu/copier"
)

type Core struct {
	config                              coreconfig.CoreConfig
	acquire                             acquire.Acquire
	trigger                             trigger.Trigger
	images                              imgcache.ImgCache
	request                             uirequest.UIRequest
	chanTrigger                         chan bool
	chanJobStop, chanJobStopped         chan bool
	chanWorkertStop, chanWorkertStopped chan bool
	configLoaded                        bool
	active                              bool
}

var mut sync.RWMutex

func (co *Core) Reset() error {

	mut.Lock()
	defer mut.Unlock()

	co.configLoaded = false
	co.config.SetUUID("")

	return nil
}

func (co *Core) LoadConfig(conf coreconfig.CoreConfig) error {

	mut.Lock()
	defer mut.Unlock()

	if co.configLoaded {
		return nil
	}

	co.config.Lock()
	copier.Copy(&co.config, &conf)
	co.config.Unlock()

	if webcam.IsDeviceAvailable(co.config.Acquire.Device.GetSource()) == false {
		co.config.SetEnabled(false)
		return nil
	}

	err := co.request.LoadConfig(co.config.UUID)
	if err != nil {
		return err
	}

	err = co.acquire.TestWebcam(co.config.Acquire)
	if err != nil {
		return err
	}

	co.setPreview()

	co.configLoaded = true

	if co.config.Enabled && co.config.Record {
		return co.start()
	}

	return nil
}

func (co *Core) Start(conf coreconfig.CoreConfig) error {

	mut.Lock()
	defer mut.Unlock()

	if co.active {
		return nil
	}

	co.config.Lock()
	copier.Copy(&co.config, &conf)
	co.config.Unlock()

	return co.start()
}

func (co *Core) start() error {

	if co.active {
		return nil
	}

	if co.config.Enabled == false {
		return errors.New("device is disabled")
	}

	err := co.acquire.Start(co.config.Acquire)
	if err != nil {
		return err
	}

	err = functions.CreateDir(co.config.GetFileRoot())
	if err != nil {
		return err
	}

	co.active = true

	co.chanTrigger = make(chan bool)
	co.chanJobStop = make(chan bool)
	co.chanJobStopped = make(chan bool)
	co.chanWorkertStop = make(chan bool)
	co.chanWorkertStopped = make(chan bool)

	go co.jobGetImage()
	go co.job()

	if co.config.Trigger.GetSource() == triggerconfig.MOTION {
		time.Sleep(co.config.GetSetup())
	}

	return co.trigger.Start(co.config.Trigger, &co.images, co.chanTrigger)
}

func (co *Core) Stop() error {

	mut.Lock()
	defer mut.Unlock()

	return co.stop()
}

func (co *Core) stop() error {

	if co.active == false {
		return nil
	}

	err := co.trigger.Stop()
	if err != nil {
		return err
	}

	co.chanWorkertStop <- true
	<-co.chanWorkertStopped
	co.chanJobStop <- true
	<-co.chanJobStopped

	close(co.chanTrigger)
	close(co.chanJobStop)
	close(co.chanJobStopped)
	close(co.chanWorkertStop)
	close(co.chanWorkertStopped)

	co.images.Clear()

	co.setPreview()

	co.active = false

	return co.acquire.Stop()
}

func (co *Core) Restart(conf coreconfig.CoreConfig) error {

	mut.Lock()
	defer mut.Unlock()

	if co.active == false {
		return nil
	}

	err := co.stop()
	if err != nil {
		return err
	}

	co.config.Lock()
	copier.Copy(&co.config, &conf)
	co.config.Unlock()

	return co.start()
}

func (co *Core) Exit() error {

	mut.Lock()
	defer mut.Unlock() //func (co *Core) setPreview(res acquireconfig.Resolution) {

	return co.stop()
}

func (co *Core) SetUUID(uid string) {

	mut.Lock()
	defer mut.Unlock()

	co.config.SetUUID(uid)
}

func (co *Core) GetUUID() string {

	mut.RLock()
	defer mut.RUnlock()

	return co.config.GetUUID()
}

func (co *Core) ListConfigs() ([]acquireconfig.DeviceConfig, error) {

	return co.acquire.ListConfigs()
}

func (co *Core) SetPreview(res acquireconfig.Resolution) {

	mut.Lock()
	defer mut.Unlock()

	co.setPreviewRes(res)
}

func (co *Core) setPreview() {

	var res acquireconfig.Resolution

	if co.config.Acquire.Process.Width != 0 && co.config.Acquire.Process.Height != 0 {
		res = co.config.Acquire.Process.Resolution
	} else {
		res = co.config.Acquire.Device.Resolution
	}

	co.setPreviewRes(res)
}

func (co *Core) setPreviewRes(res acquireconfig.Resolution) {

	if res.Width == 0 {
		res.Width = 100
	}

	if res.Height == 0 {
		res.Height = 100
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{res.Width, res.Height}

	var img typemd.MoImage
	img.Captured = time.Now()
	img.Img = image.NewRGBA(image.Rectangle{upLeft, lowRight})

	co.request.SetLatestImage(img)
}

func (co *Core) GetNewToken(uid string) string {

	mut.Lock()
	defer mut.Unlock()

	return co.request.GetNewToken(uid)
}

func (co *Core) IsActive() bool {

	mut.RLock()
	defer mut.RUnlock()

	return co.active
}

func (co *Core) job() {

	for {

		select {
		case <-co.chanTrigger:

			console.Output("trigger", "MonMotion")

			if co.config.GetRecord() {
				go co.record(co.config.GetOverrun())
			}

		case <-co.chanJobStop:
			co.chanJobStopped <- true
			return
		}
	}
}

func (co *Core) jobGetImage() {

	var (
		err error
		img typemd.MoImage
	)

	for {

		select {
		case <-co.chanWorkertStop:
			co.chanWorkertStopped <- true
			return

		default:
			func() {

				err = nil

				for i := 0; i <= 1; i++ {

					switch i {
					case 0:
						img, err = co.acquire.GetImage()

					case 1:
						co.images.SetMaxAge(co.config.GetRecodingDuration() + co.config.Trigger.GetTimeOut())
						co.images.AddImage(img)
						co.request.SetLatestImage(img)
					}

					if err != nil {
						return
					}
				}
			}()
		}
	}
}

func (co *Core) record(delay time.Duration) {

	console.Output("record", "MonMotion")

	t := time.Now()

	time.Sleep(delay)

	err := co.images.SaveImages(co.config.GetFileRoot()+t.Format("20060102_150405.0000")+"/", co.config.GetName())
	if err != nil {
		console.Log(err, "MonMotion")
	}
}
