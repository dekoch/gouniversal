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
		Lang                 lang.DeviceAcquire
		UUID                 template.HTML
		Token                template.HTML
		Device               template.HTML
		ChbRecordChecked     template.HTMLAttr
		CmbPreRecodingPeriod template.HTML
		CmbOverrunPeriod     template.HTML
		CmbSetupPeriod       template.HTML
		CmbDeviceConfig      template.HTML
		CmbProcessResolution template.HTML
		ChbCropChecked       template.HTMLAttr
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

		for i := 0; i <= 9; i++ {

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
				}

			case 4:
				c.CmbDeviceConfig, err = cmbDeviceConfig(config, dev)

			case 5:
				c.CmbProcessResolution, err = cmbProcessResolution(config)

			case 6:
				c.CmbPreRecodingPeriod, err = cmbPreRecodingPeriod(config)

			case 7:
				c.CmbOverrunPeriod, err = cmbOverrunPeriod(config)

			case 8:
				c.CmbSetupPeriod, err = cmbSetupPeriod(config)

			case 9:
				c.Token = template.HTML(dev.GetNewToken(nav.User.UUID))
				c.Device = template.HTML(id)

				if config.GetRecord() {
					c.ChbRecordChecked = template.HTMLAttr("checked")
				} else {
					c.ChbRecordChecked = template.HTMLAttr("")
				}

				if config.Acquire.Process.GetCrop() {
					c.ChbCropChecked = template.HTMLAttr("checked")
				} else {
					c.ChbCropChecked = template.HTMLAttr("")
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

func cmbProcessResolution(config *coreconfig.CoreConfig) (template.HTML, error) {

	opt := config.Acquire.GetProcessConfig()
	tag := "<select name=\"processresolution\">"
	// none
	tag += "<option value=\"-1\""
	if opt.Height == 0 &&
		opt.Width == 0 {

		tag += " selected"
	}
	tag += "></option>"
	// 640×480, 800×600, 960×720, 1024×768

	var cfgs [4]acquireconfig.Resolution
	cfgs[0].Width = 640
	cfgs[0].Height = 480
	cfgs[1].Width = 800
	cfgs[1].Height = 600
	cfgs[2].Width = 960
	cfgs[2].Height = 720
	cfgs[3].Width = 1024
	cfgs[3].Height = 768

	for _, cfg := range cfgs {

		tag += "<option value=\"" + strconv.Itoa(cfg.Width) + "x" + strconv.Itoa(cfg.Height) + "\""
		if opt.Height == cfg.Height &&
			opt.Width == cfg.Width {

			tag += " selected"
		}
		tag += ">" + strconv.Itoa(cfg.Width) + "x" + strconv.Itoa(cfg.Height) + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func cmbPreRecodingPeriod(config *coreconfig.CoreConfig) (template.HTML, error) {

	opt := config.GetPreRecoding().Milliseconds()
	tag := "<select name=\"prerecodingperiod\">"

	for i := 0; i <= 10000; i += 500 {

		tag += "<option value=\"" + strconv.Itoa(i) + "\""
		if opt == int64(i) {
			tag += " selected"
		}
		tag += ">" + strconv.Itoa(i) + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func cmbOverrunPeriod(config *coreconfig.CoreConfig) (template.HTML, error) {

	opt := config.GetOverrun().Milliseconds()
	tag := "<select name=\"overrunperiod\">"

	for i := 0; i <= 10000; i += 500 {

		tag += "<option value=\"" + strconv.Itoa(i) + "\""
		if opt == int64(i) {
			tag += " selected"
		}
		tag += ">" + strconv.Itoa(i) + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func cmbSetupPeriod(config *coreconfig.CoreConfig) (template.HTML, error) {

	opt := config.GetSetup().Milliseconds()
	tag := "<select name=\"setupperiod\">"

	for i := 0; i <= 10000; i += 500 {

		tag += "<option value=\"" + strconv.Itoa(i) + "\""
		if opt == int64(i) {
			tag += " selected"
		}
		tag += ">" + strconv.Itoa(i) + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func edit(config *coreconfig.CoreConfig, dev *core.Core, r *http.Request) error {

	var (
		err                  error
		selPreRecoding       string
		intPreRecoding       int
		selOverrun           string
		intOverrun           int
		selSetup             string
		intSetup             int
		selDeviceConfig      string
		noDeviceConfig       int
		selProcessResolution string
		dcfgs                []acquireconfig.DeviceConfig
		dcfg                 acquireconfig.DeviceConfig
		pcfg                 acquireconfig.ProcessConfig
	)

	func() {

		for i := 0; i <= 15; i++ {

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
				selProcessResolution, err = functions.CheckFormInput("processresolution", r)

			case 5:
				if functions.IsEmpty(selPreRecoding) ||
					functions.IsEmpty(selOverrun) ||
					functions.IsEmpty(selSetup) ||
					functions.IsEmpty(selDeviceConfig) ||
					functions.IsEmpty(selProcessResolution) {

					err = errors.New("bad input")
				}

			case 6:
				intPreRecoding, err = strconv.Atoi(selPreRecoding)

			case 7:
				intOverrun, err = strconv.Atoi(selOverrun)

			case 8:
				intSetup, err = strconv.Atoi(selSetup)

			case 9:
				if r.FormValue("chbrecord") == "checked" {
					config.SetRecord(true)
				} else {
					config.SetRecord(false)
				}

				config.SetPreRecoding(intPreRecoding)
				config.SetOverrun(intOverrun)
				config.SetSetup(intSetup)

			case 10:
				dcfgs, err = dev.ListConfigs()

			case 11:
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

			case 12:
				// 640×480, 800×600, 960×720, 1024×768
				switch selProcessResolution {
				case "-1":
					pcfg.Width = 0
					pcfg.Height = 0

				case "640x480":
					pcfg.Width = 640
					pcfg.Height = 480

				case "800x600":
					pcfg.Width = 800
					pcfg.Height = 600

				case "960x720":
					pcfg.Width = 960
					pcfg.Height = 720

				case "1024x768":
					pcfg.Width = 1024
					pcfg.Height = 768
				}

				err = config.Acquire.SetProcessConfig(pcfg)

			case 13:
				if r.FormValue("chbcrop") == "checked" {
					config.Acquire.Process.SetCrop(true)
				} else {
					config.Acquire.Process.SetCrop(false)
				}

			case 14:
				config.Acquire.Device.SetResolution(dcfg.Resolution)
				config.Acquire.Device.SetFPS(dcfg.FPS)

			case 15:
				if pcfg.Width != 0 && pcfg.Height != 0 {
					dev.SetPreview(pcfg.Resolution)
				} else {
					dev.SetPreview(dcfg.Resolution)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
