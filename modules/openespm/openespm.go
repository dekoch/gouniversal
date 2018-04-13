package openespm

import (
	"gouniversal/modules/openespm/appManagement"
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/request"
	"gouniversal/modules/openespm/ui"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

	globalOESPM.UiConfig.AppFileRoot = "data/ui/openespm/1.0/"

	globalOESPM.AppConfig.Mut.Lock()
	globalOESPM.AppConfig.File = appManagement.LoadConfig()
	globalOESPM.AppConfig.Mut.Unlock()

	globalOESPM.DeviceConfig.Mut.Lock()
	globalOESPM.DeviceConfig.File = deviceManagement.LoadConfig()
	globalOESPM.DeviceConfig.Mut.Unlock()

	globalOESPM.Lang.Mut.Lock()
	globalOESPM.Lang.File = langOESPM.LoadLangFiles()
	globalOESPM.Lang.Mut.Unlock()

	request.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	request.Exit()
}
