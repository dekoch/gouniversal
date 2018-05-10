package lang

import (
	"gouniversal/shared/config"
)

type Home struct {
	Title   string
	Name    string
	Size    string
	Options string
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

	l.Header = config.BuildHeader("en", "LangFileshare", 1.0, "Language File")

	l.Home.Title = "Fileshare"
	l.Home.Name = "Name"
	l.Home.Size = "Size"
	l.Home.Options = "Options"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
