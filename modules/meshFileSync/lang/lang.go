package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Search struct {
	Menu        string
	Title       string
	SelectAll   string
	Name        string
	Size        string
	Date        string
	Download    string
	ClearList   string
	UpdateFiles string
	NewFiles    string
}

type Transfers struct {
	Menu             string
	Title            string
	SelectAll        string
	Name             string
	Size             string
	Date             string
	Abort            string
	PendingDownloads string
	PendingUploads   string
}

type Settings struct {
	Menu            string
	Title           string
	Apply           string
	MaxFileSize     string
	AutoAddFiles    string
	AutoUpdateFiles string
	AutoDeleteFiles string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header    config.FileHeader
	Search    Search
	Transfers Transfers
	Settings  Settings
	Alert     Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangMeshFS", 1.0, "Language File")

	l.Search.Menu = "Mesh"
	l.Search.Title = "Search"
	l.Search.SelectAll = "Select All"
	l.Search.Name = "Name"
	l.Search.Size = "Size"
	l.Search.Date = "Date"
	l.Search.Download = "Download Selected Files"
	l.Search.ClearList = "Clear List"
	l.Search.UpdateFiles = "Files to Update"
	l.Search.NewFiles = "New Files"

	l.Transfers.Menu = "Mesh"
	l.Transfers.Title = "Transfers"
	l.Transfers.SelectAll = "Select All"
	l.Transfers.Name = "Name"
	l.Transfers.Size = "Size"
	l.Transfers.Date = "Date"
	l.Transfers.Abort = "Abort Selected Files"
	l.Transfers.PendingDownloads = "Pending Downloads"
	l.Transfers.PendingUploads = "Pending Uploads"

	l.Settings.Menu = "Mesh"
	l.Settings.Title = "Settings"
	l.Settings.Apply = "Apply"
	l.Settings.MaxFileSize = "Max File Size"
	l.Settings.AutoAddFiles = "Auto Add Files"
	l.Settings.AutoUpdateFiles = "Auto Update Files"
	l.Settings.AutoDeleteFiles = "Auto Delete Files"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
