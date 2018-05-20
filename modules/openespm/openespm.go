package openespm

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/openespm/globalOESPM"
	"github.com/dekoch/gouniversal/modules/openespm/langOESPM"
	"github.com/dekoch/gouniversal/modules/openespm/request"
	"github.com/dekoch/gouniversal/modules/openespm/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
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
