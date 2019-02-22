package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Home struct {
	Menu           string
	Title          string
	ShowLastDay    string
	ShowLast7Days  string
	ShowLast30Days string
}

type StationList struct {
	Menu       string
	Title      string
	AddStation string
	Name       string
	Company    string
	Street     string
	City       string
	Options    string
	Edit       string
}

type StationEdit struct {
	Title   string
	Apply   string
	Delete  string
	UUID    string
	Name    string
	Company string
	Street  string
	City    string
	URL     string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header      config.FileHeader
	Home        Home
	StationList StationList
	StationEdit StationEdit
	Alert       Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangGasPrice", 1.0, "Language File")

	l.Home.Menu = "GasPrice"
	l.Home.Title = "Home"
	l.Home.ShowLastDay = "Last 24h"
	l.Home.ShowLast7Days = "Last 7d"
	l.Home.ShowLast30Days = "Last 30d"

	l.StationList.Menu = "GasPrice"
	l.StationList.Title = "List"
	l.StationList.AddStation = "Add Station"
	l.StationList.Name = "Name"
	l.StationList.Company = "Company"
	l.StationList.Street = "Street"
	l.StationList.City = "City"
	l.StationList.Options = "Options"
	l.StationList.Edit = "Edit"

	l.StationEdit.Title = "Edit"
	l.StationEdit.Apply = "Apply"
	l.StationEdit.Delete = "Delete"
	l.StationEdit.UUID = "UUID"
	l.StationEdit.Name = "Name"
	l.StationEdit.Company = "Company"
	l.StationEdit.Street = "Street"
	l.StationEdit.City = "City"
	l.StationEdit.URL = "Url"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
