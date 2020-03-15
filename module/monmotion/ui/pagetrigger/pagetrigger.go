package pagetrigger

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/core"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/core/trigger/triggerconfig"
	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/module/monmotion/ui/menu"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/s7conn"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	for _, dev := range global.Config.GetDevices() {

		if dev.Enabled {
			nav.Sitemap.Register("", "App:MonMotion:Trigger$UUID="+dev.UUID, page.Lang.Device.Title+" "+dev.Name)
		}
	}
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	menu.Render("trigger", page, nav, r)

	type Content struct {
		Lang                 lang.DeviceTrigger
		CmbSource            template.HTML
		ChbTriggerAfterEvent template.HTMLAttr
		CardTitle            template.HTML
		IntervalHidden       template.HTMLAttr
		Delay                template.HTML
		MotionHidden         template.HTMLAttr
		PLCHidden            template.HTMLAttr
		PLCAddress           template.HTML
		PLCRack              template.HTML
		PLCSlot              template.HTML
		PLCVariable          template.HTML
	}
	var c Content

	c.Lang = page.Lang.Device.DeviceTrigger

	func() {

		var (
			err    error
			dev    *core.Core
			config *coreconfig.CoreConfig
		)

		// Form input
		id := nav.Parameter("UUID")

		for i := 0; i <= 7; i++ {

			switch i {
			case 0:
				dev, err = global.GetCore(id)

			case 1:
				config, err = global.Config.GetDevice(id)

			case 2:
				switch r.FormValue("editsource") {
				case "apply":
					err = editSource(config, r)
					if err == nil {
						err = global.Config.SaveConfig()
					}

					if err == nil {
						err = dev.Restart(*config)
					}
				}

			case 3:
				switch r.FormValue("editinterval") {
				case "apply":
					err = editInterval(config, r)
					if err == nil {
						err = global.Config.SaveConfig()
					}

					if err == nil {
						err = dev.Restart(*config)
					}
				}

			case 4:
				switch r.FormValue("editplc") {
				case "apply":
					err = editPLC(config, r)
					if err == nil {
						err = global.Config.SaveConfig()
					}

					if err == nil {
						err = dev.Restart(*config)
					}

				case "test":
					err := testPLC(config, page, nav)
					if err != nil {
						alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
					}
				}

			case 5:
				switch r.FormValue("editmotion") {
				case "apply":
					err = editMotion(id, r)
					if err == nil {
						err = global.Config.SaveConfig()
					}

					if err == nil {
						err = dev.Restart(*config)
					}
				}

			case 6:
				c.CmbSource, err = cmbSource(config, c.Lang)

				if config.Trigger.GetTriggerAfterEvent() {
					c.ChbTriggerAfterEvent = template.HTMLAttr("checked")
				} else {
					c.ChbTriggerAfterEvent = template.HTMLAttr("")
				}

			case 7:
				c.CardTitle = template.HTML(config.Trigger.GetSource())

				c.IntervalHidden = template.HTMLAttr("hidden")
				c.PLCHidden = template.HTMLAttr("hidden")
				c.MotionHidden = template.HTMLAttr("hidden")

				switch config.Trigger.GetSource() {
				case triggerconfig.INTERVAL:
					c.IntervalHidden = template.HTMLAttr("")

					cfg := config.Trigger.GetIntervalConfig()
					c.Delay = template.HTML(strconv.Itoa(cfg.Delay))

				case triggerconfig.PLC:
					c.PLCHidden = template.HTMLAttr("")

					cfg := config.Trigger.GetPLCConfig()
					c.PLCAddress = template.HTML(cfg.Address)
					c.PLCRack = template.HTML(strconv.Itoa(cfg.Rack))
					c.PLCSlot = template.HTML(strconv.Itoa(cfg.Slot))
					c.PLCVariable = template.HTML(cfg.Variable)

				case triggerconfig.MOTION:
					c.MotionHidden = template.HTMLAttr("")
				}
			}

			if err != nil {
				alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				return
			}
		}
	}()

	p, err := functions.PageToString(global.Config.UIFileRoot+"trigger.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func cmbSource(config *coreconfig.CoreConfig, lang lang.DeviceTrigger) (template.HTML, error) {

	opt := config.Trigger.GetSource()
	tag := "<select name=\"source\">"

	tag += "<option value=\"disabled\""
	if opt == triggerconfig.DISABLED {
		tag += " selected"
	}
	tag += ">" + lang.Disabled + "</option>"

	tag += "<option value=\"interval\""
	if opt == triggerconfig.INTERVAL {
		tag += " selected"
	}
	tag += ">" + lang.Interval + "</option>"

	tag += "<option value=\"plc\""
	if opt == triggerconfig.PLC {
		tag += " selected"
	}
	tag += ">" + lang.PLC + "</option>"

	tag += "<option value=\"motion\""
	if opt == triggerconfig.MOTION {
		tag += " selected"
	}
	tag += ">" + lang.Motion + "</option>"

	tag += "</select>"

	return template.HTML(tag), nil
}

func editSource(config *coreconfig.CoreConfig, r *http.Request) error {

	var err error

	func() {

		var (
			selSource string
			src       triggerconfig.Source
		)

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				selSource, err = functions.CheckFormInput("source", r)

			case 1:
				if functions.IsEmpty(selSource) {
					err = errors.New("bad input")
				}

			case 2:
				switch selSource {
				case "disabled":
					src = triggerconfig.DISABLED

				case "interval":
					src = triggerconfig.INTERVAL

				case "plc":
					src = triggerconfig.PLC

				case "motion":
					src = triggerconfig.MOTION

				default:
					err = errors.New("invalid trigger source")
				}

			case 3:
				err = config.Trigger.SetSource(src)

			case 4:
				if r.FormValue("chbtriggerafterevent") == "checked" {
					config.Trigger.SetTrigger(true)
				} else {
					config.Trigger.SetTrigger(false)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func editInterval(config *coreconfig.CoreConfig, r *http.Request) error {

	var err error

	func() {

		var (
			strDelay string
			intDelay int
		)

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				strDelay, err = functions.CheckFormInput("delay", r)

			case 1:
				if functions.IsEmpty(strDelay) {
					err = errors.New("bad input")
				}

			case 2:
				intDelay, err = strconv.Atoi(strDelay)

			case 3:
				var cfg triggerconfig.SourceInterval
				cfg.Delay = intDelay

				err = config.Trigger.SetIntervalConfig(cfg)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func editPLC(config *coreconfig.CoreConfig, r *http.Request) error {

	var err error

	func() {

		var (
			address  string
			strRack  string
			intRack  int
			strSlot  string
			intSlot  int
			variable string
		)

		for i := 0; i <= 7; i++ {

			switch i {
			case 0:
				address, err = functions.CheckFormInput("plcaddress", r)

			case 1:
				strRack, err = functions.CheckFormInput("plcrack", r)

			case 2:
				strSlot, err = functions.CheckFormInput("plcslot", r)

			case 3:
				variable, err = functions.CheckFormInput("plcvariable", r)

			case 4:
				if functions.IsEmpty(address) ||
					functions.IsEmpty(strRack) ||
					functions.IsEmpty(strSlot) ||
					functions.IsEmpty(variable) {
					err = errors.New("bad input")
				}

			case 5:
				intRack, err = strconv.Atoi(strRack)

			case 6:
				intSlot, err = strconv.Atoi(strSlot)

			case 7:
				var cfg triggerconfig.SourcePLC
				cfg.Address = address
				cfg.Rack = intRack
				cfg.Slot = intSlot
				cfg.Variable = variable

				err = config.Trigger.SetPLCConfig(cfg)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func editMotion(id string, r *http.Request) error {

	var err error

	func() {

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func testPLC(config *coreconfig.CoreConfig, page *typemd.Page, nav *navigation.Navigation) error {

	var err error

	func() {

		var (
			cfg   triggerconfig.SourcePLC
			plc   s7conn.S7Conn
			conn  *s7conn.Connection
			val   interface{}
			state bool
		)

		for i := 0; i <= 6; i++ {

			switch i {
			case 0:
				cfg = config.Trigger.GetPLCConfig()

			case 1:
				err = plc.AddPLC(cfg.Address, cfg.Rack, cfg.Slot, 1, 200*time.Millisecond, 600*time.Millisecond)

			case 2:
				conn, err = plc.GetConnection(cfg.Address)

			case 3:
				defer conn.Release()

			case 4:
				val, err = conn.Client.Read(cfg.Variable)

			case 5:
				switch val.(type) {
				case bool:
					state = val.(bool)

				default:
					err = errors.New("unsupported variable " + cfg.Variable)
				}

			case 6:
				if state {
					alert.Message(alert.SUCCESS, page.Lang.Alert.Success, cfg.Variable+"=true", page.Lang.Device.DeviceTrigger.TestConnection, nav.User.UUID)
				} else {
					alert.Message(alert.SUCCESS, page.Lang.Alert.Success, cfg.Variable+"=false", page.Lang.Device.DeviceTrigger.TestConnection, nav.User.UUID)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
