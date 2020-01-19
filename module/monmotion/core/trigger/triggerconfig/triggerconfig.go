package triggerconfig

import (
	"errors"
	"sync"
	"time"
)

type Source string

const (
	MOTION   Source = "Motion"
	PLC      Source = "PLC"
	INTERVAL Source = "Interval"
	DISABLED Source = "disabled"
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

type SourceInterval struct {
	Delay int // second
}

type TriggerConfig struct {
	TriggerAfterEvent bool
	CheckIntvl        int // millisecond
	Source            Source
	Motion            SourceMotion
	PLC               SourcePLC
	Interval          SourceInterval
}

var mut sync.RWMutex

func (hc *TriggerConfig) LoadDefaults() {

	hc.TriggerAfterEvent = false
	hc.Source = DISABLED
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

	hc.Interval.Delay = 10
}

func (hc *TriggerConfig) Lock() {

	mut.Lock()
}

func (hc *TriggerConfig) Unlock() {

	mut.Unlock()
}

func (hc *TriggerConfig) SetSource(source Source) error {

	mut.Lock()
	defer mut.Unlock()

	if source == DISABLED ||
		source == PLC ||
		source == MOTION ||
		source == INTERVAL {

		hc.Source = source
		return nil
	}

	return errors.New("invalid trigger source")
}

func (hc *TriggerConfig) GetSource() Source {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Source
}

func (hc *TriggerConfig) SetTrigger(afterevent bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.TriggerAfterEvent = afterevent
}

func (hc *TriggerConfig) GetTriggerAfterEvent() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.TriggerAfterEvent
}

func (hc *TriggerConfig) GetCheckIntvl() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.CheckIntvl) * time.Millisecond
}

func (hc *TriggerConfig) SetMotionConfig(motion SourceMotion) {

	mut.Lock()
	defer mut.Unlock()

	hc.Motion = motion
}

func (hc *TriggerConfig) GetMotionConfig() SourceMotion {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Motion
}

func (hc *TriggerConfig) SetPLCConfig(plc SourcePLC) error {

	mut.Lock()
	defer mut.Unlock()

	hc.PLC = plc

	return nil
}

func (hc *TriggerConfig) GetPLCConfig() SourcePLC {

	mut.RLock()
	defer mut.RUnlock()

	return hc.PLC
}

func (hc *TriggerConfig) SetIntervalConfig(interval SourceInterval) error {

	mut.Lock()
	defer mut.Unlock()

	hc.Interval = interval

	return nil
}

func (hc *TriggerConfig) GetIntervalConfig() SourceInterval {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Interval
}

func (hc *TriggerConfig) GetTimeOut() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	if hc.Source == MOTION {
		return time.Duration(hc.Motion.TimeOut) * time.Millisecond
	}

	return time.Duration(0)
}

func (hc *SourceMotion) GetTimeSpan() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.TimeSpan) * time.Millisecond
}
