package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu         string
	Title        string
	Time         string
	State        string
	Step         string
	CurrentPos   string
	NextWaypoint string
	Waypoint     string
	Start        string
	Stop         string
	Upload       string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header config.FileHeader
	Home   Home
	Alert  Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangNav", 1.0, "Language File")

	l.Home.Menu = "Tools"
	l.Home.Title = "GPSNav"
	l.Home.Time = "Time"
	l.Home.State = "State"
	l.Home.Step = "Step"
	l.Home.CurrentPos = "Current Position"
	l.Home.NextWaypoint = "Next Waypoint"
	l.Home.Waypoint = "Waypoint"
	l.Home.Start = "Start"
	l.Home.Stop = "Stop"
	l.Home.Upload = "Upload"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
