package pagedevicelist

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dekoch/gouniversal/module/monmotion/core"
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/webcam"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.DeviceList.Menu, "App:MonMotion:DeviceList", page.Lang.DeviceList.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang       lang.DeviceList
		DeviceList template.HTML
	}
	var c Content

	c.Lang = page.Lang.DeviceList

	var (
		err      error
		redirect bool
	)

	func() {

		availableDevs := webcam.FindDevices()

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				var newDev bool
				newDev, err = scanNewDevices(availableDevs)
				if newDev {
					redirect = true
				}

			case 1:
				switch r.FormValue("edit") {
				case "apply":
					err = edit(r)
					if err == nil {
						err = global.Config.SaveConfig()
						redirect = true
					}
				}

			case 2:
				c.DeviceList, err = tblDeviceList(availableDevs)
			}

			if err != nil {
				alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				return
			}

			if redirect {
				return
			}
		}
	}()

	if err == nil {
		if redirect {
			nav.RedirectPath("App:MonMotion:DeviceList", false)
			return
		}
	}

	p, err := functions.PageToString(global.Config.UIFileRoot+"devicelist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func tblDeviceList(availabledevs []string) (template.HTML, error) {

	tbody := ""

	for _, dev := range global.Config.GetDevices() {

		tbody += "<tr>"
		tbody += "<td><input type=\"checkbox\" name=\"selecteddevice\" value=\"" + dev.UUID + "\""

		if dev.Enabled {
			tbody += " checked"
		}

		missing := true

		for _, availableDev := range availabledevs {

			if dev.Acquire.Device.Source == availableDev {
				missing = false
			}
		}

		if missing {
			tbody += " disabled"
		}

		tbody += "></td>"

		tbody += "<td>" + dev.Name + "</td>"
		tbody += "<td>" + dev.Acquire.Device.Source + "</td>"
		tbody += "</tr>"
	}

	return template.HTML(tbody), nil
}

func edit(r *http.Request) error {

	var err error

	func() {

		enabledDevs := global.Config.GetEnabledDevices()

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				global.Config.SetEnabledDevices(r.Form["selecteddevice"])

			case 1:
				// stop disabled devices
				for _, disabledDev := range global.Config.GetDisabledDevices() {

					for _, enabledDev := range enabledDevs {

						if disabledDev == enabledDev {
							err = global.FreeCore(disabledDev)
							if err != nil {
								return
							}
						}
					}
				}

			case 2:
				// init new devices
				for _, enabledDev := range global.Config.GetEnabledDevices() {

					if global.IsCoreAvailable(enabledDev) == false {

						err = startCore(enabledDev)
						if err != nil {
							cfg, err := global.Config.GetDevice(enabledDev)
							if err != nil {
								return
							}

							cfg.SetEnabled(false)

							return
						}
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func startCore(uid string) error {

	var err error

	func() {

		var (
			dev *core.Core
			cfg *coreconfig.CoreConfig
		)

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dev, err = global.GetFreeCore()

			case 1:
				cfg, err = global.Config.GetDevice(uid)

			case 2:
				err = dev.LoadConfig(*cfg)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func scanNewDevices(availabledevs []string) (bool, error) {

	foundNew := false

	for _, dev := range availabledevs {

		var cfg coreconfig.CoreConfig
		cfg.LoadDefaults()
		cfg.Name = strings.Replace(dev, "/dev/", "", -1)
		cfg.FileRoot += cfg.Name + "/"
		cfg.Acquire.Device.Source = dev

		if global.Config.AddNewDevice(cfg) {
			foundNew = true
		}
	}

	return foundNew, nil
}
