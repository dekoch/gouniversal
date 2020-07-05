package s7backup

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/core"
	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	err := global.Config.LoadConfig()
	if err != nil {
		console.Log(err, "")
	}

	ui.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	err = core.LoadConfig()
	if err != nil {
		console.Log(err, "")
	}
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit(em *types.ExitMessage) {

}
