package langOESPM

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const LangDir = "data/lang/openespm/"

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
	Title   string
	Name    string
	App     string
	State   string
	States  SettingsAppState
	Comment string
	Apply   string
	Delete  string
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

type File struct {
	Header           config.FileHeader
	Settings         Settings
	Alert            Alert
	SimpleSwitchV1x0 SimpleSwitchV1x0
}

type Global struct {
	Mut  sync.Mutex
	File []File
}

func SaveLang(lang File, n string) error {

	lang.Header = config.BuildHeader(n, "LangOpenESPM", 1.0, "Language File")

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		lang.Settings.Title = "Settings"

		lang.Settings.App.Title = "openESPM Apps"
		lang.Settings.App.List.Title = "App List"
		lang.Settings.App.List.AddApp = "Add App"
		lang.Settings.App.List.Name = "Name"
		lang.Settings.App.List.App = "App"
		lang.Settings.App.List.Options = "Options"
		lang.Settings.App.List.Edit = "Edit"

		lang.Settings.App.Edit.Title = "App Edit"
		lang.Settings.App.Edit.Name = "Name"
		lang.Settings.App.Edit.App = "App"
		lang.Settings.App.Edit.State = "State"
		lang.Settings.App.Edit.States.Active = "Active"
		lang.Settings.App.Edit.States.Inactive = "Inactive"
		lang.Settings.App.Edit.Comment = "Comment"
		lang.Settings.App.Edit.Apply = "Apply"
		lang.Settings.App.Edit.Delete = "Delete"

		lang.Settings.Device.Title = "openESPM Devices"
		lang.Settings.Device.List.Title = "Device List"
		lang.Settings.Device.List.AddDevice = "Add Device"
		lang.Settings.Device.List.Name = "Name"
		lang.Settings.Device.List.App = "App"
		lang.Settings.Device.List.Options = "Options"
		lang.Settings.Device.List.Edit = "Edit"

		lang.Settings.Device.Edit.Title = "Device Edit"
		lang.Settings.Device.Edit.Name = "Name"
		lang.Settings.Device.Edit.App = "App"
		lang.Settings.Device.Edit.State = "State"
		lang.Settings.Device.Edit.States.Active = "Active"
		lang.Settings.Device.Edit.States.Inactive = "Inactive"
		lang.Settings.Device.Edit.Comment = "Comment"
		lang.Settings.Device.Edit.Apply = "Apply"
		lang.Settings.Device.Edit.Delete = "Delete"

		lang.Alert.Success = "Success"
		lang.Alert.Info = "Info"
		lang.Alert.Warning = "Warning"
		lang.Alert.Error = "Error"

		lang.SimpleSwitchV1x0.On = "On"
		lang.SimpleSwitchV1x0.Off = "Off"
		lang.SimpleSwitchV1x0.DeviceSettings = "Device Settings"
	}

	b, err := json.Marshal(lang)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(LangDir+n, b)

	return err
}

func LoadLang(n string) File {

	var lg File

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		SaveLang(lg, n)
	}

	f := new(file.File)
	b, err := f.ReadFile(LangDir + n)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &lg)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(lg.Header, "LangOpenESPM") == false {
		log.Fatal("wrong config")
	}

	return lg
}

func LoadLangFiles() []File {

	lg := make([]File, 0)

	if _, err := os.Stat(LangDir + "en"); os.IsNotExist(err) {
		// if not found, create default file
		var newlg File
		SaveLang(newlg, "en")
	}

	files, err := ioutil.ReadDir(LangDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fl := range files {

		var langfile File

		f := new(file.File)
		b, err := f.ReadFile(LangDir + fl.Name())
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(b, &langfile)
		if err != nil {
			log.Fatal(err)
		}

		if config.CheckHeader(langfile.Header, "LangOpenESPM") {

			lg = append(lg, langfile)
		}

	}

	return lg
}
