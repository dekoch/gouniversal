package typenav

import (
	"time"

	"github.com/dekoch/gouniversal/module/gpsnav/lang"
)

type StepType string

const (
	StepInit          StepType = "Init"
	StepWaitForGPSFix StepType = "WaitForGPSFix"
	StepGPSFix        StepType = "GPSFix"
	StepCheckGPS      StepType = "CheckGPS"
	StepNavigate      StepType = "Navigate"
	StepTrack         StepType = "Track"
	StepEnd           StepType = "End"
)

type Pos struct {
	Time    time.Time // Time of fix.
	Lat     float64   // Latitude.
	Lon     float64   // Longitude.
	Fix     string    // Quality of fix.
	Sat     int64     // Number of satellites in use.
	HDOP    float64   // Horizontal dilution of precision.
	Ele     float64   // Altitude.
	Name    string
	Comment string
	Valid   bool // Pos is valid
}

type Page struct {
	Content string
	Lang    lang.LangFile
}
