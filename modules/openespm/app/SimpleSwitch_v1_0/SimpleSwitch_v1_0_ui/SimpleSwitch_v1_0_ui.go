package SimpleSwitch_v1_0_ui

import (
	"fmt"
	"gouniversal/modules/openespm/app/SimpleSwitch_v1_0"
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/alert"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"net/http"
)

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	var config SimpleSwitch_v1_0.AppConfig
	var sw string

	type app struct {
		Lang   langOESPM.SimpleSwitch_v1_0
		Switch template.HTML
	}
	var a app

	a.Lang = page.Lang.SimpleSwitch_v1_0

	// load device
	var dev typesOESPM.Device

	globalOESPM.DeviceConfig.Mut.Lock()
	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search selcted device
		if nav.IsNext(globalOESPM.DeviceConfig.File.Devices[u].UUID) {

			dev = globalOESPM.DeviceConfig.File.Devices[u]
		}
	}
	globalOESPM.DeviceConfig.Mut.Unlock()

	if dev.UUID == "" {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, "device not found", nav.CurrentPath, nav.User.UUID)
		return
	}

	// init new device
	if functions.IsEmpty(dev.Config) {
		dev.Config = SimpleSwitch_v1_0.InitConfig()
	}

	// read config
	dev.Unmarshal(&config)

	// toggle switch
	newState, err := functions.CheckFormInput("switch", r)
	if err == nil {
		if functions.IsEmpty(newState) == false {
			if newState == "on" {
				config.Switch = true
			} else {
				config.Switch = false
			}
		}
	}

	// save device
	err = dev.Marshal(config)
	if err == nil {
		deviceManagement.SaveDevice(dev.UUID, dev)
	}

	// render page
	if config.Switch {
		sw = "<button class=\"btn btn-success fa fa-toggle-on\" type=\"submit\" name=\"switch\" value=\"off\" title=\"" + a.Lang.On + "\"></button>"
	} else {
		sw = "<button class=\"btn btn-danger fa fa-toggle-off\" type=\"submit\" name=\"switch\" value=\"on\" title=\"" + a.Lang.Off + "\"></button>"
	}
	a.Switch = template.HTML(sw)

	templ, err := template.ParseFiles(globalOESPM.UiConfig.AppFileRoot + "SimpleSwitch_v1_0/app.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, a)
}
