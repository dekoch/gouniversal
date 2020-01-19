package trigger

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core/imgcache"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/analyse"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/s7conn"
	"github.com/dekoch/gouniversal/shared/timeout"

	"github.com/jinzhu/copier"
)

type MotionState int

const (
	READY MotionState = 1 + iota
	ACTIVATED
)

type Trigger struct {
	config  triggerconfig.TriggerConfig
	images  *imgcache.ImgCache
	analyse analyse.Analyse

	chanTrigger                           chan bool
	chanTriggerStart, chanTriggerFinished chan bool
	chanJobStop, chanJobStopped           chan bool
	chanWorkerStop, chanWorkerStopped     chan bool
	active                                bool
}

var mut sync.RWMutex

func (tr *Trigger) Start(conf triggerconfig.TriggerConfig, images *imgcache.ImgCache, trigger chan bool) error {

	mut.Lock()
	defer mut.Unlock()

	if tr.active {
		return nil
	}

	tr.config.Lock()
	copier.Copy(&tr.config, &conf)
	tr.config.Unlock()

	if tr.config.Source == triggerconfig.DISABLED {
		return nil
	}

	tr.images = images
	tr.chanTrigger = trigger

	tr.chanTriggerStart = make(chan bool)
	tr.chanTriggerFinished = make(chan bool)

	tr.chanWorkerStop = make(chan bool)
	tr.chanWorkerStopped = make(chan bool)

	tr.chanJobStop = make(chan bool)
	tr.chanJobStopped = make(chan bool)

	go tr.job()

	switch tr.config.Source {
	case triggerconfig.MOTION:
		tr.analyse.LoadConfig()

		go tr.jobDetectMotion()

	case triggerconfig.PLC:

		go tr.jobPLC()

	case triggerconfig.INTERVAL:

		go tr.jobInterval()

	default:
		return errors.New("invalid source")
	}

	tr.active = true

	return nil
}

func (tr *Trigger) Stop() error {

	mut.Lock()
	defer mut.Unlock()

	if tr.active == false {
		return nil
	}

	if tr.config.Source == triggerconfig.DISABLED {
		return nil
	}

	tr.chanJobStop <- true
	<-tr.chanJobStopped

	tr.chanWorkerStop <- true
	<-tr.chanWorkerStopped

	close(tr.chanJobStop)
	close(tr.chanJobStopped)

	close(tr.chanWorkerStop)
	close(tr.chanWorkerStopped)

	close(tr.chanTriggerStart)
	close(tr.chanTriggerFinished)

	tr.active = false

	return nil
}

func (tr *Trigger) job() {

	timerTrigger := time.NewTimer(tr.config.GetCheckIntvl())

	var toTrigger timeout.TimeOut
	toTrigger.Start(999)

	for {

		select {
		case <-timerTrigger.C:
			timerTrigger.Stop()

			toTrigger.Reset()

			select {
			case <-tr.chanJobStop:
				tr.chanJobStopped <- true
				return

			default:
				tr.chanTriggerStart <- true
			}

		case <-tr.chanTriggerFinished:
			pause := tr.config.GetCheckIntvl() - (time.Duration(toTrigger.ElapsedMillis()) * time.Millisecond)

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
		newImg, oldImg mdimg.MDImage
		res            analyse.Result
	)

	motionState = READY

	for {

		select {
		case <-tr.chanTriggerStart:
			func() {

				err = nil

				config = tr.config.GetMotionConfig()

				if config.AutoTune {

					config.AutoTune = false
					config.Threshold = 0
					tr.config.SetMotionConfig(config)

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
						res, err = tr.analyse.AnalyseImage(oldImg, newImg, config.Threshold)

					case 4:
						if res.Threshold > config.Threshold {

							config.Threshold = res.Threshold

							tr.config.SetMotionConfig(config)
						}

					case 5:
						if config.Threshold > 0 {

							if tr.config.GetTriggerAfterEvent() {

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

		case <-tr.chanWorkerStop:
			tr.chanWorkerStopped <- true
			return
		}
	}
}

func (tr *Trigger) jobPLC() {

	var (
		err                error
		plc                s7conn.S7Conn
		config             triggerconfig.SourcePLC
		conn               *s7conn.Connection
		newValue, oldValue bool
		val                interface{}
	)

	for {

		select {
		case <-tr.chanTriggerStart:
			func() {

				err = nil

				config = tr.config.GetPLCConfig()

				for i := 0; i <= 5; i++ {

					switch i {
					case 0:
						err = plc.AddPLC(config.Address, config.Rack, config.Slot, 1, 200*time.Millisecond, tr.config.GetCheckIntvl()*3)

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

							if tr.config.GetTriggerAfterEvent() {

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

		case <-tr.chanWorkerStop:
			tr.chanWorkerStopped <- true
			return
		}
	}
}

func (tr *Trigger) jobInterval() {

	var to timeout.TimeOut
	to.Start(tr.config.GetIntervalConfig().Delay * 1000)

	for {

		select {
		case <-tr.chanTriggerStart:

			to.SetTimeOut(tr.config.GetIntervalConfig().Delay * 1000)

			if to.Elapsed() {

				to.Reset()
				tr.chanTrigger <- true
			}

			tr.chanTriggerFinished <- true

		case <-tr.chanWorkerStop:
			tr.chanWorkerStopped <- true
			return
		}
	}
}
