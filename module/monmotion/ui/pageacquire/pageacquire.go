package pageacquire

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/monmotion/core"
	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/acquireconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/module/monmotion/ui/menu"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	for _, dev := range global.Config.GetDevices() {

		if dev.Enabled {
			nav.Sitemap.Register(page.Lang.Device.Menu, "App:MonMotion:Acquire$UUID="+dev.UUID, page.Lang.Device.Title+" "+dev.Name)
		}
	}
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	menu.Render("acquire", page, nav, r)

	type Content struct {
		Lang              lang.DeviceAcquire
		UUID              template.HTML
		Token             template.HTML
		Device            template.HTML
		ChbRecordChecked  template.HTMLAttr
		PreRecodingPeriod template.HTML
		OverrunPeriod     template.HTML
		SetupPeriod       template.HTML
		CmbDeviceConfig   template.HTML
	}
	var c Content

	c.Lang = page.Lang.Device.DeviceAcquire

	c.UUID = template.HTML(nav.User.UUID)

	func() {

		var (
			err    error
			dev    *core.Core
			config *coreconfig.CoreConfig
		)

		// Form input
		id := nav.Parameter("UUID")

		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				dev, err = global.GetCore(id)

			case 1:
				config, err = global.Config.GetDevice(id)

			case 2:
				switch r.FormValue("edit") {
				case "apply":
					err = edit(config, dev, r)
					if err == nil {
						err = global.Config.SaveConfig()
					}

					if err == nil {
						err = dev.Restart(*config)
					}
				}

			case 3:
				switch r.FormValue("control") {
				case "start":
					err = dev.Start(*config)

				case "stop":
					err = dev.Stop()

				case "trigger":
					err = dev.ManualTrigger()
				}

			case 4:
				c.CmbDeviceConfig, err = cmbDeviceConfig(config, dev)

			case 5:
				c.PreRecodingPeriod = template.HTML(strconv.FormatFloat(config.GetPreRecoding().Seconds(), 'f', 0, 64))

			case 6:
				c.OverrunPeriod = template.HTML(strconv.FormatFloat(config.GetOverrun().Seconds(), 'f', 0, 64))

			case 7:
				c.SetupPeriod = template.HTML(strconv.FormatFloat(config.GetSetup().Seconds(), 'f', 0, 64))

			case 8:
				c.Token = template.HTML(dev.GetNewToken(nav.User.UUID))
				c.Device = template.HTML(id)

				if config.GetRecord() {
					c.ChbRecordChecked = template.HTMLAttr("checked")
				} else {
					c.ChbRecordChecked = template.HTMLAttr("")
				}
			}

			if err != nil {
				alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				return
			}
		}
	}()

	p, err := functions.PageToString(global.Config.UIFileRoot+"acquire.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func cmbDeviceConfig(config *coreconfig.CoreConfig, dev *core.Core) (template.HTML, error) {

	opt := config.Acquire.GetDeviceConfig()
	cfgs, err := dev.ListConfigs()
	if err != nil {
		return template.HTML(""), err
	}

	tag := "<select name=\"deviceconfig\">"
	// auto
	tag += "<option value=\"-1\""
	if opt.Height == 0 &&
		opt.Width == 0 {

		tag += " selected"
	}
	tag += ">auto</option>"
	// device configs
	for i, cfg := range cfgs {

		tag += "<option value=\"" + strconv.Itoa(i) + "\""
		if opt.Height == cfg.Height &&
			opt.Width == cfg.Width &&
			opt.FPS == cfg.FPS {

			tag += " selected"
		}
		tag += ">" + strconv.Itoa(cfg.Width) + "x" + strconv.Itoa(cfg.Height) + " @ " + strconv.Itoa(int(cfg.FPS)) + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func edit(config *coreconfig.CoreConfig, dev *core.Core, r *http.Request) error {

	var err error

	func() {

		var (
			selPreRecoding  string
			intPreRecoding  int
			selOverrun      string
			intOverrun      int
			selSetup        string
			intSetup        int
			selDeviceConfig string
			noDeviceConfig  int
			dcfgs           []acquireconfig.DeviceConfig
			dcfg            acquireconfig.DeviceConfig
		)

		for i := 0; i <= 12; i++ {

			switch i {
			case 0:
				selPreRecoding, err = functions.CheckFormInput("prerecodingperiod", r)

			case 1:
				selOverrun, err = functions.CheckFormInput("overrunperiod", r)

			case 2:
				selSetup, err = functions.CheckFormInput("setupperiod", r)

			case 3:
				selDeviceConfig = r.FormValue("deviceconfig")

			case 4:
				if functions.IsEmpty(selPreRecoding) ||
					functions.IsEmpty(selOverrun) ||
					functions.IsEmpty(selSetup) ||
					functions.IsEmpty(selDeviceConfig) {

					err = errors.New("bad input")
				}

			case 5:
				intPreRecoding, err = strconv.Atoi(selPreRecoding)

			case 6:
				intOverrun, err = strconv.Atoi(selOverrun)

			case 7:
				intSetup, err = strconv.Atoi(selSetup)

			case 8:
				if r.FormValue("chbrecord") == "checked" {
					config.SetRecord(true)
				} else {
					config.SetRecord(false)
				}

				config.SetPreRecoding(intPreRecoding)
				config.SetOverrun(intOverrun)
				config.SetSetup(intSetup)

			case 9:
				dcfgs, err = dev.ListConfigs()

			case 10:
				noDeviceConfig, err = strconv.Atoi(selDeviceConfig)

				switch noDeviceConfig {
				case -1:
					dcfg.Width = 0
					dcfg.Height = 0
					dcfg.FPS = 0

				default:
					if noDeviceConfig >= 0 && noDeviceConfig < len(dcfgs) {
						dcfg = dcfgs[noDeviceConfig]
					}
				}

			case 11:
				config.Acquire.Device.SetResolution(dcfg.Resolution)
				config.Acquire.Device.SetFPS(dcfg.FPS)

			case 12:
				dev.SetPreview(dcfg.Resolution)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
