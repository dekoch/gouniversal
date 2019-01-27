package core

import (
	"errors"
	"sync"
	"time"

	"github.com/adrianmo/go-nmea"

	"github.com/dekoch/gouniversal/module/gpsnav/geo"
	"github.com/dekoch/gouniversal/module/gpsnav/global"
	"github.com/dekoch/gouniversal/module/gpsnav/gps"
	"github.com/dekoch/gouniversal/module/gpsnav/route"
	"github.com/dekoch/gouniversal/module/gpsnav/tracker"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/sbool"
	"github.com/dekoch/gouniversal/shared/sint"
	"github.com/dekoch/gouniversal/shared/timeout"
)

var (
	mut          sync.RWMutex
	delay        time.Duration
	step         typenav.StepType
	mygps        gps.Gps
	stop         sbool.Sbool
	state        sint.Sint
	watchdog     timeout.TimeOut
	lastWaypoint sint.Sint
	bearing      geo.Geo
)

func LoadConfig() {

	err := mygps.OpenPort(global.Config.GetGPSPort(), global.Config.GetGPSBaud(), global.Config.GetGPSTimeOut())
	if err != nil {
		console.Log(err, "")
		return
	}

	mygps.LoadConfig()

	state.Set(typenav.STOPPED)

	go job()
}

func Exit() {

	mygps.ClosePort()
}

func Start(startWaypoint int) error {

	if state.Get() == typenav.STOPPED {

		stop.UnSet()
		state.Set(typenav.RUNNING)

		watchdog.SetTimeOut(10000)
		watchdog.Enable(true)
		watchdog.Reset()

		go run(startWaypoint)

		return nil
	}

	return errors.New("nav already running")
}

func Stop() error {

	if state.Get() == typenav.RUNNING {

		stop.Set()

		return nil
	}

	return errors.New("nav not running")
}

func GetState() int {

	return state.Get()
}

func GetBearing() (float64, error) {

	if state.Get() != typenav.RUNNING {
		return 0.0, errors.New("nav not running")
	}

	return bearing.GetBearing()
}

func job() {

	tCheckWatchdog := time.NewTicker(time.Duration(2) * time.Second)

	for {
		select {
		case <-tCheckWatchdog.C:
			if state.Get() == typenav.RUNNING &&
				watchdog.Elapsed() {

				console.Log("error: watchog", "nav")
				state.Set(typenav.STOPPED)
				Start(lastWaypoint.Get())
			}
		}
	}
}

func run(startWaypoint int) {

	console.Log("started", "nav")

	var (
		err    error
		tr     tracker.Tracker
		rt     route.Route
		curPos typenav.Pos
	)

	rt.Init(global.Config.GetRouteFilePath())
	err = rt.ReadFile()
	if err != nil {
		console.Log(err, "")
		return
	}

	rt.StartRoute()
	err = rt.SetWaypoint(startWaypoint)
	if err != nil {
		console.Log(err, "")
		return
	}

	lastWaypoint.Set(startWaypoint)
	nextwpt, _ := rt.GetNextWaypoint()
	console.Log(nextwpt.Name, "nav")

	t := time.Now()
	var fileName string
	fileName += t.Format("20060102")
	fileName += "_"
	fileName += t.Format("150405")
	fileName += ".gpx"
	tr.Init(global.Config.GetGPXRoot()+fileName, global.Config.GetGPXTrackDist(), true)
	tr.Enable(true)

	nextStep(typenav.StepInit, 0, true)
	//oldStep := typenav.StepInit

	for {

		err = nil

		func() {

			curPos = mygps.GetPos()

			switch getStep() {
			case typenav.StepInit:

				mygps.ResetTimeOut()

				nextStep(typenav.StepWaitForGPSFix, global.Config.GetGPSTimeOut()*2, true)

			case typenav.StepWaitForGPSFix:

				if mygps.GetTimeOut() {
					return
				}

				fix := curPos.Fix

				if fix != nmea.GPS {
					return
				}

				nextStep(typenav.StepGPSFix, 0, true)

			case typenav.StepGPSFix:

				bearing.SetStartPos(curPos)

				nextStep(typenav.StepCheckGPS, 0, true)

			case typenav.StepCheckGPS:

				if mygps.GetTimeOut() {
					err = errors.New("GPS TimeOut")
					return
				}

				if curPos.Fix != nmea.GPS {
					err = errors.New("lost GPS fix")
					return
				}

				nextStep(typenav.StepNavigate, 0, false)

			case typenav.StepNavigate:

				nextStep(typenav.StepTrack, 0, false)

				onwpt, err := rt.OnWaypoint(curPos, global.Config.GetRouteWptMaxDist())
				if err != nil {
					return
				}

				if onwpt {
					// if we reached waypoint, save current pos
					wpt, err := rt.GetNextWaypoint()
					if err == nil {

						pos := curPos
						pos.Name = wpt.Name

						go func(p typenav.Pos) {
							tr.NewWaypoint(p)
							tr.WriteFile()
						}(pos)
					}

					rt.NextWaypoint()
					lastWaypoint.Set(rt.GetWaypointNo())

					bearing.SetStartPos(curPos)

					nextwpt, _ := rt.GetNextWaypoint()
					console.Log(nextwpt.Name, "nav")
				}
				// calculate bearing
				nextwpt, err := rt.GetNextWaypoint()
				if err != nil {
					return
				}

				bearing.SetTargetPos(nextwpt)
				bearing.SetCurrentPos(curPos)

				//bearing, _ := bearing.GetTargetBearing()
				//fmt.Println(bearing)

			case typenav.StepTrack:

				nextStep(typenav.StepEnd, 0, false)

				go func(p typenav.Pos) {
					tr.CheckPos(p)
				}(curPos)

			case typenav.StepEnd:
				nextStep(typenav.StepCheckGPS, 200, false)

			default:
				nextStep(typenav.StepInit, 0, true)
			}
		}()

		if err != nil {

			console.Log("error: step "+string(getStep())+" "+err.Error(), "nav")

			nextStep(typenav.StepInit, 0, true)
		}

		/*if getStep() != oldStep {
			oldStep = getStep()

			fmt.Printf("step %v\n", oldStep)
		}*/

		if stop.IsSet() {

			console.Log("stopped", "nav")
			state.Set(typenav.STOPPED)
			return
		}

		watchdog.Reset()

		time.Sleep(getDelay())
	}
}

func nextStep(s typenav.StepType, m int, log bool) {

	mut.Lock()
	defer mut.Unlock()

	step = s
	delay = time.Duration(m) * time.Millisecond

	if log {
		console.Log("step "+string(step), "nav")
	}
}

func getStep() typenav.StepType {

	mut.RLock()
	defer mut.RUnlock()

	return step
}

func setDelay(m int) {

	mut.Lock()
	defer mut.Unlock()

	delay = time.Duration(m) * time.Millisecond
}

func getDelay() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return delay
}
