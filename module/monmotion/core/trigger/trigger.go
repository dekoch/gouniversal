package trigger

import (
	"errors"
	"fmt"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/imgcache"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/analyse"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/s7conn"
	"github.com/dekoch/gouniversal/shared/timeout"
)

type MotionState int

const (
	READY MotionState = 1 + iota
	ACTIVATED
)

type Trigger struct {
	triggerconfig.TriggerConfig
	images  *imgcache.ImgCache
	analyse analyse.Analyse

	chanTriggerStart    chan bool
	chanTriggerFinished chan bool

	chanTrigger chan bool
}

func (tr *Trigger) LoadConfig(images *imgcache.ImgCache, trigger chan bool) error {

	tr.images = images
	tr.chanTrigger = trigger

	tr.chanTriggerStart = make(chan bool)
	tr.chanTriggerFinished = make(chan bool)

	switch tr.GetSource() {
	case triggerconfig.MOTION:
		tr.analyse.LoadConfig()

		go tr.jobDetectMotion()

	case triggerconfig.PLC:

		go tr.jobPLC()

	default:
		return errors.New("invalid source")
	}

	go tr.job()

	return nil
}

func (tr *Trigger) job() {

	timerTrigger := time.NewTimer(tr.GetCheckIntvl())

	var toTrigger timeout.TimeOut
	toTrigger.Start(999)

	for {

		select {
		case <-timerTrigger.C:
			timerTrigger.Stop()

			toTrigger.Reset()

			tr.chanTriggerStart <- true

		case <-tr.chanTriggerFinished:
			pause := tr.GetCheckIntvl() - (time.Duration(toTrigger.ElapsedMillis()) * time.Millisecond)

			//fmt.Println(pause)

			timerTrigger.Reset(pause)
		}
	}
}

func (tr *Trigger) jobDetectMotion() {

	var (
		err            error
		config         triggerconfig.SourceMotion
		motionState    MotionState
		to             timeout.TimeOut
		newImg, oldImg typemd.MoImage
		res            analyse.Result
	)

	motionState = READY

	for {

		<-tr.chanTriggerStart

		func() {

			err = nil

			config = tr.GetMotionConfig()

			if config.AutoTune {

				config.AutoTune = false
				config.Threshold = 0
				tr.SetMotionConfig(config)

				tr.analyse.EnableAutoTune(config.TuneTime, config.TuneStep)
			}

			for i := 0; i <= 5; i++ {

				switch i {
				case 0:
					if tr.images.GetImageCnt() < 2 {
						return
					}

				case 1:
					oldImg, err = tr.images.GetOldImage(time.Duration(config.TimeSpan) * time.Millisecond)

				case 2:
					newImg, err = tr.images.GetLatestImage()

				case 3:
					res, err = tr.analyse.AnalyseImage(&oldImg, &newImg, config.Threshold)

				case 4:
					if res.Threshold > config.Threshold {

						config.Threshold = res.Threshold

						tr.SetMotionConfig(config)
					}

				case 5:
					if config.Threshold > 0 {

						if tr.GetTriggerAfterEvent() {

							if res.Px > 0 {

								if motionState != ACTIVATED {
									console.Output("moving", "MonMotion")
								}

								motionState = ACTIVATED
								to.Start(config.TimeOut)
							} else {

								if motionState == ACTIVATED && to.Elapsed() {

									motionState = READY
									tr.chanTrigger <- true
								}
							}
						} else {

							if res.Px > 0 {

								if motionState == READY {

									console.Output("moving", "MonMotion")

									motionState = ACTIVATED
									tr.chanTrigger <- true
								}

								to.Start(config.TimeOut)

							} else {

								if to.Elapsed() {
									motionState = READY
								}
							}
						}
					}
				}

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()

		tr.chanTriggerFinished <- true
	}
}

func (tr *Trigger) jobPLC() {

	var (
		err                error
		config             triggerconfig.SourcePLC
		plc                s7conn.S7Conn
		conn               *s7conn.Connection
		newValue, oldValue bool
		val                interface{}
	)

	for {

		<-tr.chanTriggerStart

		func() {

			err = nil

			config = tr.GetPLCConfig()

			for i := 0; i <= 5; i++ {

				switch i {
				case 0:
					err = plc.AddPLC(config.Address, config.Rack, config.Slot, 1, 200*time.Millisecond, tr.GetCheckIntvl()*3)

				case 1:
					conn, err = plc.GetConnection(config.Address)

				case 2:
					defer conn.Release()

				case 3:
					val, err = conn.Client.Read(config.Variable)

				case 4:
					switch val.(type) {
					case bool:
						newValue = val.(bool)

					default:
						err = errors.New("unsupported variable " + config.Variable)
					}

				case 5:
					if newValue != oldValue {

						oldValue = newValue

						if tr.GetTriggerAfterEvent() {

							if newValue == false {
								tr.chanTrigger <- true
							}
						} else {

							if newValue {
								tr.chanTrigger <- true
							}
						}
					}
				}

				if err != nil {
					console.Output(err, "MonMotion")
					return
				}
			}
		}()

		tr.chanTriggerFinished <- true
	}
}
