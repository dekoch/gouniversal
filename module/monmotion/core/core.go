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
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/sbool"

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
	triggerState                        sbool.Sbool
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

	err = co.images.LoadConfig(co.config.UUID)
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

	co.images.SetRAMSettings(co.config.Trigger.Motion.GetTimeSpan())
	co.images.SetDBSettings(co.config.GetRecodingDuration()+co.config.Trigger.GetTimeOut(), 60)

	err := co.acquire.Start(co.config.Acquire)
	if err != nil {
		return err
	}

	err = functions.CreateDir(co.config.GetFileRoot())
	if err != nil {
		return err
	}

	co.active = true

	co.triggerState.UnSet()

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

	co.setPreview()

	co.active = false

	err = co.acquire.Stop()
	if err != nil {
		return err
	}

	return co.images.Exit()
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
	defer mut.Unlock()

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

func (co *Core) SetPreview(res typemd.Resolution) {

	mut.Lock()
	defer mut.Unlock()

	co.setPreviewRes(res)
}

func (co *Core) setPreview() {

	co.setPreviewRes(co.config.Acquire.Device.Resolution)
}

func (co *Core) setPreviewRes(res typemd.Resolution) {

	if res.Width == 0 {
		res.Width = 100
	}

	if res.Height == 0 {
		res.Height = 100
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{res.Width, res.Height}

	var preview mdimg.MDImage
	preview.Captured = time.Now()
	preview.Resolution = res

	err := preview.EncodeImage(image.NewRGBA(image.Rectangle{upLeft, lowRight}))
	if err != nil {
		return
	}

	co.request.SetLatestImage(preview)
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

			co.triggerState.Set()

		case <-co.chanJobStop:
			co.chanJobStopped <- true
			return
		}
	}
}

func (co *Core) jobGetImage() {

	var (
		err error
		img mdimg.MDImage
	)

	recordEnabled := co.config.GetRecord()
	triggerEnabled := co.config.Trigger.GetSource() != triggerconfig.DISABLED

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
						if recordEnabled {
							img.Trigger = co.triggerState.IsSet()

							var record bool

							if triggerEnabled {
								if img.Trigger {

									co.triggerState.UnSet()

									img.PreRecoding = co.config.GetPreRecoding().Seconds()
									img.Overrun = co.config.GetOverrun().Seconds()

									record = true
								}
							}

							if record {
								go co.record(co.config.GetOverrun())
							}

							co.images.AddImage(img, true)
						} else {
							co.images.AddImage(img, false)
						}

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

	trigger := time.Now()

	time.Sleep(delay)

	err := co.images.SaveImages(trigger, trigger.Add(delay))
	if err != nil {
		console.Log(err, "MonMotion")
	}
}
