package console

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/console/global"
	"github.com/dekoch/gouniversal/modules/console/lang"
	"github.com/dekoch/gouniversal/modules/console/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	ui.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

}
