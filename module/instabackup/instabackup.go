package instabackup

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/instabackup/core"
	"github.com/dekoch/gouniversal/module/instabackup/global"
	"github.com/dekoch/gouniversal/module/instabackup/lang"
	"github.com/dekoch/gouniversal/module/instabackup/request"
	"github.com/dekoch/gouniversal/module/instabackup/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	global.Tokens.SetMaxTokens(global.Config.GetMaxTokens())

	core.LoadConfig()
	request.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	err := global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "")
	}
}
