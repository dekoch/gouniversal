package TempHumV1x0ui

import (
	"gouniversal/modules/openespm/app/TempHumV1x0"
	"gouniversal/modules/openespm/deviceConfig"
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

	edit, _ := functions.CheckFormInput("edit", r)

	var app TempHumV1x0.AppConfig
	var dev TempHumV1x0.DeviceConfig
	var device deviceConfig.Device

	type content struct {
		Lang      langOESPM.TempHumV1x0
		SelectDev template.HTML
		Switch    template.HTML
	}
	var c content

	c.Lang = page.Lang.TempHumV1x0

	var err error
	err = nil
	for i := 0; i <= 9; i++ {
		if err == nil {
			switch i {
			case 0:
				// init new app
				if functions.IsEmpty(page.App.Config) {
					page.App.Config = TempHumV1x0.InitAppConfig()
				}

			case 1:
				// app config to struct
				err = page.App.Unmarshal(&app)

			case 2:
				// form input
				// selected device
				selDevice, err := functions.CheckFormInput("device", r)
				if err == nil {
					if functions.IsEmpty(selDevice) == false {
						app.DeviceUUID = selDevice
					}
				}

			case 3:
				// load device config
				if app.DeviceUUID != "" {
					device, err = globalOESPM.DeviceConfig.Get(app.DeviceUUID)
				}
				// init new device
				if functions.IsEmpty(device.Config) {
					device.Config = TempHumV1x0.InitDeviceConfig()
				}

			case 4:
				// device config to struct
				err = device.Unmarshal(&dev)

			case 5:
				// form input
				// toggle switch
				/*newState, err := functions.CheckFormInput("switch", r)
				if err == nil {
					if functions.IsEmpty(newState) == false {
						if newState == "on" {
							dev.Switch = true
						} else {
							dev.Switch = false
						}
					}
				}*/

			case 6:
				// struct to app config
				err = page.App.Marshal(app)

			case 7:
				// struct to device config
				err = device.Marshal(dev)

			case 8:
				// save device to ram
				err = globalOESPM.DeviceConfig.Edit(app.DeviceUUID, device)

			case 9:
				// save device to file
				if edit == "apply" {
					err = globalOESPM.DeviceConfig.SaveConfig()
				}
			}
		}
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err.Error(), nav.CurrentPath, nav.User.UUID)
	}

	// render page

	c.SelectDev = deviceManagement.HTMLSelectDevice("device", page.App.App, app.DeviceUUID)

	p, err := functions.PageToString(globalOESPM.UiConfig.AppFileRoot+"TempHumV1x0/app.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
