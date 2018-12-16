package fileshare

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/fileshare/global"
	"github.com/dekoch/gouniversal/module/fileshare/lang"
	"github.com/dekoch/gouniversal/module/fileshare/request"
	"github.com/dekoch/gouniversal/module/fileshare/ui"
	"github.com/dekoch/gouniversal/module/fileshare/upload"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
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
