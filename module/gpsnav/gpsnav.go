package gpsnav

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/gpsnav/core"
	"github.com/dekoch/gouniversal/module/gpsnav/global"
	"github.com/dekoch/gouniversal/module/gpsnav/lang"
	"github.com/dekoch/gouniversal/module/gpsnav/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/gpsnav/", en, "en")

	global.Config.LoadConfig()
	core.LoadConfig()

	ui.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	core.Exit()
	global.Config.SaveConfig()
}
