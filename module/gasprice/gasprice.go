package gasprice

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/gasprice/core"
	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/lang"
	"github.com/dekoch/gouniversal/module/gasprice/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	ui.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	core.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	core.Exit()
}
