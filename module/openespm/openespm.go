package openespm

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/module/openespm/request"
	"github.com/dekoch/gouniversal/module/openespm/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.UiConfig.AppFileRoot = "data/ui/openespm/1.0/"

	global.AppConfig.LoadConfig()
	global.DeviceConfig.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/openespm/", en, "en")

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
