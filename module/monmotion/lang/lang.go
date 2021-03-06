package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type DeviceList struct {
	Menu   string
	Title  string
	Enable string
	Name   string
	Source string
	Apply  string
}

type DeviceMenu struct {
	Acquire string
	Trigger string
}

type DeviceAcquire struct {
	Control            string
	Live               string
	Start              string
	Stop               string
	Trigger            string
	Settings           string
	Record             string
	KeepAllSequences   string
	PreRecodingPeriod  string
	OverrunPeriod      string
	SetupPeriod        string
	InputResolutionFPS string
	OutputResolution   string
	CropOutput         string
	Apply              string
	Console            string
}

type DeviceTrigger struct {
	Source            string
	TriggerAfterEvent string
	Apply             string
	Disabled          string
	Interval          string
	Delay             string
	Motion            string
	PLC               string
	Address           string
	Rack              string
	Slot              string
	Variable          string
	TestConnection    string
}

type Device struct {
	Menu          string
	Title         string
	DeviceMenu    DeviceMenu
	DeviceAcquire DeviceAcquire
	DeviceTrigger DeviceTrigger
}

type Viewer struct {
	Menu    string
	Title   string
	Refresh string
	View    string
	Export  string
	Delete  string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header     config.FileHeader
	DeviceList DeviceList
	Device     Device
	Viewer     Viewer
	Alert      Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangInstaBackup", 1.0, "Language File")

	l.DeviceList.Menu = "MonMotion"
	l.DeviceList.Title = "Device List"
	l.DeviceList.Enable = "Enable"
	l.DeviceList.Name = "Name"
	l.DeviceList.Source = "Source"
	l.DeviceList.Apply = "Apply"

	l.Device.Menu = l.DeviceList.Menu
	l.Device.Title = "Device"

	l.Device.DeviceMenu.Acquire = "Acquire"
	l.Device.DeviceMenu.Trigger = "Trigger"

	l.Device.DeviceAcquire.Control = "Control"
	l.Device.DeviceAcquire.Live = "Live"
	l.Device.DeviceAcquire.Start = "Start"
	l.Device.DeviceAcquire.Stop = "Stop"
	l.Device.DeviceAcquire.Trigger = "Trigger"
	l.Device.DeviceAcquire.Settings = "Settings"
	l.Device.DeviceAcquire.Record = "Record"
	l.Device.DeviceAcquire.KeepAllSequences = "Keep all Sequences"
	l.Device.DeviceAcquire.PreRecodingPeriod = "Pre Recoding Period"
	l.Device.DeviceAcquire.OverrunPeriod = "Overrun Period"
	l.Device.DeviceAcquire.SetupPeriod = "Setup Period"
	l.Device.DeviceAcquire.InputResolutionFPS = "Input Resolution/FPS"
	l.Device.DeviceAcquire.OutputResolution = "Output Resolution"
	l.Device.DeviceAcquire.CropOutput = "Crop Output"
	l.Device.DeviceAcquire.Apply = "Apply"
	l.Device.DeviceAcquire.Console = "Console"

	l.Device.DeviceTrigger.Source = "Source"
	l.Device.DeviceTrigger.TriggerAfterEvent = "Trigger After Event"
	l.Device.DeviceTrigger.Apply = "Apply"
	l.Device.DeviceTrigger.Disabled = "disabled"
	l.Device.DeviceTrigger.Interval = "Interval"
	l.Device.DeviceTrigger.Delay = "Delay"
	l.Device.DeviceTrigger.Motion = "Motion"
	l.Device.DeviceTrigger.PLC = "PLC"
	l.Device.DeviceTrigger.Address = "Address"
	l.Device.DeviceTrigger.Rack = "Rack"
	l.Device.DeviceTrigger.Slot = "Slot"
	l.Device.DeviceTrigger.Variable = "Variable (DBX/M/I/O)"
	l.Device.DeviceTrigger.TestConnection = "Test Connection"

	l.Viewer.Menu = l.DeviceList.Menu
	l.Viewer.Title = "Viewer"
	l.Viewer.Refresh = "Refresh"
	l.Viewer.View = "View Sequence"
	l.Viewer.Export = "Export Sequence"
	l.Viewer.Delete = "Delete Sequence"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
