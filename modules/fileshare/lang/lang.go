package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu      string
	Title     string
	NewFolder string
	Name      string
	Size      string
	Options   string
	Edit      string
	Delete    string
	Upload    string
}

type Edit struct {
	Title string
	Name  string
	Apply string
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
	Edit   Edit
	Alert  Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangFileshare", 1.0, "Language File")

	l.Home.Menu = "Tools"
	l.Home.Title = "Fileshare"
	l.Home.NewFolder = "new Folder"
	l.Home.Name = "Name"
	l.Home.Size = "Size"
	l.Home.Options = "Options"
	l.Home.Edit = "Edit"
	l.Home.Delete = "Delete"
	l.Home.Upload = "Upload"

	l.Edit.Title = "Edit"
	l.Edit.Name = "Name"
	l.Edit.Apply = "Apply"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
