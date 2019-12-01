package triggerconfig

import (
	"sync"
	"time"
)

type Source string

const (
	MOTION Source = "Motion"
	PLC    Source = "PLC"
)

type SourceMotion struct {
	TimeSpan  int    // millisecond
	Threshold uint32 // px
	TimeOut   int    // millisecond
	AutoTune  bool
	TuneStep  uint32 // px
	TuneTime  int    // millisecond
}

type SourcePLC struct {
	Address  string
	Rack     int
	Slot     int
	Variable string
}

type TriggerConfig struct {
	TriggerAfterEvent bool
	CheckIntvl        int // millisecond
	Source            Source
	Motion            SourceMotion
	PLC               SourcePLC
	mut               sync.RWMutex
}

func (hc *TriggerConfig) LoadDefaults() {

	hc.TriggerAfterEvent = false
	hc.Source = MOTION
	hc.CheckIntvl = 200

	hc.Motion.TimeSpan = 1000
	hc.Motion.Threshold = 0
	hc.Motion.TimeOut = 3000
	hc.Motion.AutoTune = true
	hc.Motion.TuneTime = 10000
	hc.Motion.TuneStep = 2000

	hc.PLC.Address = "192.168.178.240"
	hc.PLC.Rack = 0
	hc.PLC.Slot = 2
	hc.PLC.Variable = "DB1000.DBX0.0"
}

func (hc *TriggerConfig) SetSource(source Source) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Source = source
}

func (hc *TriggerConfig) GetSource() Source {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.Source
}

func (hc *TriggerConfig) SetTrigger(afterevent bool) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.TriggerAfterEvent = afterevent
}

func (hc *TriggerConfig) GetTriggerAfterEvent() bool {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.TriggerAfterEvent
}

func (hc *TriggerConfig) GetCheckIntvl() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return time.Duration(hc.CheckIntvl) * time.Millisecond
}

func (hc *TriggerConfig) SetMotionConfig(motion SourceMotion) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.Motion = motion
}

func (hc *TriggerConfig) GetMotionConfig() SourceMotion {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.Motion
}

func (hc *TriggerConfig) SetPLCConfig(plc SourcePLC) {

	hc.mut.Lock()
	defer hc.mut.Unlock()

	hc.PLC = plc
}

func (hc *TriggerConfig) GetPLCConfig() SourcePLC {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	return hc.PLC
}

func (hc *TriggerConfig) GetTimeOut() time.Duration {

	hc.mut.RLock()
	defer hc.mut.RUnlock()

	if hc.Source == MOTION {
		return time.Duration(hc.Motion.TimeOut) * time.Millisecond
	}

	return time.Duration(0)
}
