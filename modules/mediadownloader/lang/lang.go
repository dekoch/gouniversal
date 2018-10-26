package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu                    string
	Title                   string
	PleaseEnterUrl          string
	Url                     string
	Find                    string
	Link                    string
	Download                string
	DownloadFinished        string
	NotSupportedContentType string
	NoFileFound             string
	SupportedFileExtensions string
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

	l.Header = config.BuildHeader("en", "LangMediaDownloader", 1.0, "Language File")

	l.Home.Menu = "Tools"

	l.Home.Title = "MediaDownloader"
	l.Home.PleaseEnterUrl = "please enter url"
	l.Home.Url = "Url"
	l.Home.Find = "Find"
	l.Home.Download = "Find+Download"
	l.Home.Link = "Link"
	l.Home.DownloadFinished = "Download Finished!"
	l.Home.NotSupportedContentType = "not supported content type"
	l.Home.NoFileFound = "no file found"
	l.Home.SupportedFileExtensions = "supported file extensions"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
