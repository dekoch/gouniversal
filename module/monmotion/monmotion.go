package monmotion

import (
	"net/http"
	"strings"

	"github.com/dekoch/gouniversal/module/monmotion/core/acquire/webcam"
	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	for _, dev := range webcam.FindDevices() {

		var cfg coreconfig.CoreConfig
		cfg.LoadDefaults()
		cfg.Name = strings.Replace(dev, "/dev/", "", -1)
		cfg.FileRoot += cfg.Name + "/"
		cfg.Acquire.Device.Source = dev

		global.Config.AddNewDevice(cfg)
	}

	global.AddCores(global.Config.GetNoCores())

	for _, cfg := range global.Config.GetDevices() {

		if webcam.IsDeviceAvailable(cfg.Acquire.Device.GetSource()) == false {

			d, err := global.Config.GetDevice(cfg.GetUUID())
			if err == nil {
				d.SetEnabled(false)
			}

			continue
		}

		if cfg.GetEnabled() == false {
			continue
		}

		dev, err := global.GetFreeCore()
		if err != nil {
			console.Log(err, "")
			continue
		}

		err = dev.LoadConfig(cfg)
		if err != nil {
			console.Log(err, "")
		}
	}

	global.Config.SaveConfig()

	ui.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit(em *types.ExitMessage) {

	err := global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "")
	}

	for _, config := range global.Config.GetDevices() {

		dev, err := global.GetCore(config.UUID)
		if err != nil {
			console.Log(err, "")
			continue
		}

		err = dev.Exit()
		if err != nil {
			console.Log(err, "")
		}
	}
}
