package lang

import (
	"gouniversal/shared/config"
)

type Home struct {
	Title     string
	NewFolder string
	Name      string
	Size      string
	Options   string
	Delete    string
	Upload    string
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
	l.Home.NewFolder = "new Folder"
	l.Home.Name = "Name"
	l.Home.Size = "Size"
	l.Home.Options = "Options"
	l.Home.Delete = "Delete"
	l.Home.Upload = "Upload"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
