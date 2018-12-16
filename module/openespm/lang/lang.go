package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type SettingsAppList struct {
	Title   string
	AddApp  string
	Name    string
	App     string
	Options string
	Edit    string
}

type SettingsAppState struct {
	Active   string
	Inactive string
}

type SettingsAppEdit struct {
	Title     string
	Name      string
	App       string
	State     string
	States    SettingsAppState
	Comment   string
	Apply     string
	Delete    string
	InfoGroup string
}

type SettingsApp struct {
	Title string
	List  SettingsAppList
	Edit  SettingsAppEdit
}

type SettingsDeviceList struct {
	Title     string
	AddDevice string
	Name      string
	App       string
	Options   string
	Edit      string
}

type SettingsDeviceState struct {
	Active   string
	Inactive string
}

type SettingsDeviceEdit struct {
	Title   string
	Name    string
	App     string
	State   string
	States  SettingsDeviceState
	Comment string
	Apply   string
	Delete  string
}

type SettingsDevice struct {
	Title string
	List  SettingsDeviceList
	Edit  SettingsDeviceEdit
}

type Settings struct {
	Title  string
	App    SettingsApp
	Device SettingsDevice
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type SimpleSwitchV1x0 struct {
	On             string
	Off            string
	DeviceSettings string
}

type TempHumV1x0 struct {
	DeviceSettings string
}

type LangFile struct {
	Header           config.FileHeader
	Settings         Settings
	Alert            Alert
	SimpleSwitchV1x0 SimpleSwitchV1x0
	TempHumV1x0      TempHumV1x0
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangOpenESPM", 1.0, "Language File")

	l.Settings.Title = "Settings"

	l.Settings.App.Title = "openESPM Apps"
	l.Settings.App.List.Title = "App List"
	l.Settings.App.List.AddApp = "Add App"
	l.Settings.App.List.Name = "Name"
	l.Settings.App.List.App = "App"
	l.Settings.App.List.Options = "Options"
	l.Settings.App.List.Edit = "Edit"

	l.Settings.App.Edit.Title = "App Edit"
	l.Settings.App.Edit.Name = "Name"
	l.Settings.App.Edit.App = "App"
	l.Settings.App.Edit.State = "State"
	l.Settings.App.Edit.States.Active = "Active"
	l.Settings.App.Edit.States.Inactive = "Inactive"
	l.Settings.App.Edit.Comment = "Comment"
	l.Settings.App.Edit.Apply = "Apply"
	l.Settings.App.Edit.Delete = "Delete"
	l.Settings.App.Edit.InfoGroup = "Don't forget to enable new page."

	l.Settings.Device.Title = "openESPM Devices"
	l.Settings.Device.List.Title = "Device List"
	l.Settings.Device.List.AddDevice = "Add Device"
	l.Settings.Device.List.Name = "Name"
	l.Settings.Device.List.App = "App"
	l.Settings.Device.List.Options = "Options"
	l.Settings.Device.List.Edit = "Edit"

	l.Settings.Device.Edit.Title = "Device Edit"
	l.Settings.Device.Edit.Name = "Name"
	l.Settings.Device.Edit.App = "App"
	l.Settings.Device.Edit.State = "State"
	l.Settings.Device.Edit.States.Active = "Active"
	l.Settings.Device.Edit.States.Inactive = "Inactive"
	l.Settings.Device.Edit.Comment = "Comment"
	l.Settings.Device.Edit.Apply = "Apply"
	l.Settings.Device.Edit.Delete = "Delete"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	l.SimpleSwitchV1x0.On = "On"
	l.SimpleSwitchV1x0.Off = "Off"
	l.SimpleSwitchV1x0.DeviceSettings = "Device Settings"

	l.TempHumV1x0.DeviceSettings = "Device Settings"

	return l
}
