package temphumv1x0ui

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/app/temphumv1x0"
	"github.com/dekoch/gouniversal/module/openespm/deviceconfig"
	"github.com/dekoch/gouniversal/module/openespm/devicemanagement"
	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	edit, _ := functions.CheckFormInput("edit", r)

	var app temphumv1x0.AppConfig
	var dev temphumv1x0.DeviceConfig
	var device deviceconfig.Device

	type content struct {
		Lang      lang.TempHumV1x0
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
					page.App.Config = temphumv1x0.InitAppConfig()
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
					device, err = global.DeviceConfig.Get(app.DeviceUUID)
				}
				// init new device
				if functions.IsEmpty(device.Config) {
					device.Config = temphumv1x0.InitDeviceConfig()
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
				err = global.DeviceConfig.Edit(app.DeviceUUID, device)

			case 9:
				// save device to file
				if edit == "apply" {
					err = global.DeviceConfig.SaveConfig()
				}
			}
		}
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err.Error(), nav.CurrentPath, nav.User.UUID)
	}

	// render page

	c.SelectDev = devicemanagement.HTMLSelectDevice("device", page.App.App, app.DeviceUUID)

	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"TempHumV1x0/app.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
