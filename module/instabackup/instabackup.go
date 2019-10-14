package instabackup

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/instabackup/core"
	"github.com/dekoch/gouniversal/module/instabackup/global"
	"github.com/dekoch/gouniversal/module/instabackup/lang"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	core.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

}

func Exit() {

}
