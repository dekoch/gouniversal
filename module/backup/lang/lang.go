package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu            string
	Title           string
	PathFile        string
	Size            string
	DownloadArchive string
	NoItemSelected  string
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

	l.Header = config.BuildHeader("en", "LangBackup", 1.0, "Language File")

	l.Home.Menu = "Program"

	l.Home.Title = "Backup"
	l.Home.PathFile = "Path/File"
	l.Home.Size = "Size"
	l.Home.DownloadArchive = "download archive"
	l.Home.NoItemSelected = "no item selected"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
