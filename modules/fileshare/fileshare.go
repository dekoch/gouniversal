package fileshare

import (
	"gouniversal/modules/fileshare/global"
	"gouniversal/modules/fileshare/lang"
	"gouniversal/modules/fileshare/request"
	"gouniversal/modules/fileshare/ui"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

	global.Lang.Files = lang.LoadLangFiles()

	global.Config.LoadConfig()

	request.LoadConfig()

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
