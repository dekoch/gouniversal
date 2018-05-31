package logviewer

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/logviewer/global"
	"github.com/dekoch/gouniversal/modules/logviewer/lang"
	"github.com/dekoch/gouniversal/modules/logviewer/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.File.LangFileRoot, en, "en")
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

}
