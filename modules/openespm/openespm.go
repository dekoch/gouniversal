package openespm

import (
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/request"
	"gouniversal/modules/openespm/ui"
	"gouniversal/shared/language"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

	globalOESPM.UiConfig.AppFileRoot = "data/ui/openespm/1.0/"

	globalOESPM.AppConfig.LoadConfig()
	globalOESPM.DeviceConfig.LoadConfig()

	en := langOESPM.DefaultEn()
	globalOESPM.Lang = language.New("data/lang/openespm/", en, "en")

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
