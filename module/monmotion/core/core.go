package core

import (
	"fmt"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/imgcache"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

type Core struct {
	images      imgcache.ImgCache
	chanTrigger chan bool

	coreconfig.CoreConfig
	Acquire acquire.Acquire
	Trigger trigger.Trigger
}

func (co *Core) LoadConfig() {

	if co.Enabled == false {
		return
	}

	co.chanTrigger = make(chan bool)

	err := co.Acquire.Start()
	if err != nil {
		console.Log(err, "")
		return
	}

	err = functions.CreateDir(co.GetFileRoot())
	if err != nil {
		console.Log(err, "")
		return
	}

	go co.jobGetImage()
	go co.job()

	time.Sleep(co.GetSetup())

	err = co.Trigger.LoadConfig(&co.images, co.chanTrigger)
	if err != nil {
		console.Log(err, "")
		return
	}
}

func (co *Core) LoadDefaults() {

	co.LoadCoreDefaults()
	co.Acquire.LoadDefaults()
	co.Trigger.LoadDefaults()
}

func (co *Core) job() {

	for {

		select {
		case <-co.chanTrigger:

			console.Output("trigger", "MonMotion")

			if co.GetRecord() {
				go co.record(co.GetOverrun())
			}
		}
	}
}

func (co *Core) jobGetImage() {

	var (
		err error
		img typemd.MoImage
	)

	for {

		func() {

			err = nil

			for i := 0; i <= 1; i++ {

				switch i {
				case 0:
					img, err = co.Acquire.GetImage()

				case 1:
					co.images.SetMaxAge(co.GetRecodingDuration() + co.Trigger.GetTimeOut())
					co.images.AddImage(img)
				}

				if err != nil {
					return
				}
			}
		}()
	}
}

func (co *Core) record(delay time.Duration) {

	console.Output("record", "MonMotion")

	t := time.Now()

	time.Sleep(delay)

	err := co.images.SaveImages(co.GetFileRoot()+t.Format("20060102_150405.0000")+"/", co.GetName())
	if err != nil {
		fmt.Println(err)
	}
}
