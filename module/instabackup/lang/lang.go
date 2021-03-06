package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu  string
	Title string
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

	l.Header = config.BuildHeader("en", "LangInstaBackup", 1.0, "Language File")

	l.Home.Menu = "Tools"

	l.Home.Title = "InstaBackup"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
