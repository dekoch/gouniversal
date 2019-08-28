package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Title    string
	Logfiles string
	Name     string
	Size     string
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

	l.Header = config.BuildHeader("en", "LangLogViewer", 1.0, "Language File")

	l.Home.Title = "LogViewer"
	l.Home.Logfiles = "Logfiles"
	l.Home.Name = "Name"
	l.Home.Size = "Size"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
