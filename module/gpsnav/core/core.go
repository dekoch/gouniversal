package core

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/adrianmo/go-nmea"

	"github.com/dekoch/gouniversal/module/gpsnav/global"
	"github.com/dekoch/gouniversal/module/gpsnav/gps"
	"github.com/dekoch/gouniversal/module/gpsnav/route"
	"github.com/dekoch/gouniversal/module/gpsnav/tracker"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/console"
)

var (
	mut   sync.RWMutex
	delay time.Duration
	step  typenav.StepType
	mygps gps.Gps
)

func LoadConfig() {

	err := mygps.OpenPort(global.Config.GetGPSPort(), global.Config.GetGPSBaud())
	if err != nil {
		console.Log(err, "")
		return
	}

	mygps.LoadConfig()

	go run()
	//go job()
}

func Exit() {

	mygps.ClosePort()
}

func job() {

	tStepLog := time.NewTicker(time.Duration(30) * time.Second)

	for {

		select {
		case <-tStepLog.C:
			console.Log("step "+string(getStep()), "core")
		}
	}
}

func run() {

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

	nextStep(typenav.StepInit, 0)
	oldStep := typenav.StepInit

	for {

		err = nil

		func() {

			curPos = mygps.GetPos()

			switch getStep() {
			case typenav.StepInit:

				nextStep(typenav.StepWaitForGPSFix, 200)

				mygps.ResetTimeOut()

			case typenav.StepWaitForGPSFix:

				if mygps.GetTimeOut() {
					return
				}

				fix := curPos.Fix

				if fix != nmea.GPS {

					if len(fix) > 0 {
						fmt.Println(fix)
					}

					return
				}

				nextStep(typenav.StepGPSFix, 0)

			case typenav.StepGPSFix:

				global.Geo.SetStartPos(curPos)

				nextStep(typenav.StepCheckGPS, 0)

			case typenav.StepCheckGPS:

				if mygps.GetTimeOut() {
					err = errors.New("GPS TimeOut")
					return
				}

				if curPos.Fix != nmea.GPS {
					err = errors.New("lost GPS fix")
					return
				}

				nextStep(typenav.StepNavigate, 0)

			case typenav.StepNavigate:

				nextStep(typenav.StepTrack, 0)

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

					global.Geo.SetStartPos(curPos)

					nextwpt, _ := rt.GetNextWaypoint()
					console.Log(nextwpt.Name, "nav")
				}
				// calculate bearing
				nextwpt, err := rt.GetNextWaypoint()
				if err != nil {
					return
				}

				global.Geo.SetTargetPos(nextwpt)
				global.Geo.SetCurrentPos(curPos)

				//bearing, err := global.Geo.GetTargetBearing()
				//fmt.Println(bearing)

			case typenav.StepTrack:

				nextStep(typenav.StepEnd, 0)

				go func(p typenav.Pos) {
					tr.CheckPos(p)
				}(curPos)

			case typenav.StepEnd:
				nextStep(typenav.StepCheckGPS, 200)

			default:
				nextStep(typenav.StepInit, 0)
			}
		}()

		if err != nil {

			console.Log("error: step "+string(getStep())+" "+err.Error(), "core")

			nextStep(typenav.StepInit, 0)
		}

		if getStep() != oldStep {
			oldStep = getStep()

			//fmt.Printf("step %v\n", oldStep)
		}

		time.Sleep(getDelay())
	}
}

func nextStep(s typenav.StepType, m int) {

	mut.Lock()
	defer mut.Unlock()

	step = s
	delay = time.Duration(m) * time.Millisecond
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
