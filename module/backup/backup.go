package backup

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/backup/global"
	"github.com/dekoch/gouniversal/module/backup/lang"
	"github.com/dekoch/gouniversal/module/backup/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit(em *types.ExitMessage) {

	err := global.Config.Exit(em)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "")
	}
}
