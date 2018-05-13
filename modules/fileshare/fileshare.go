package fileshare

import (
	"gouniversal/modules/fileshare/global"
	"gouniversal/modules/fileshare/lang"
	"gouniversal/modules/fileshare/request"
	"gouniversal/modules/fileshare/ui"
	"gouniversal/modules/fileshare/upload"
	"gouniversal/shared/language"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/fileshare/", en, "en")

	global.Config.LoadConfig()

	request.LoadConfig()
	upload.LoadConfig()

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
