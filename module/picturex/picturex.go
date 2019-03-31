package picturex

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/picturex/global"
	"github.com/dekoch/gouniversal/module/picturex/lang"
	"github.com/dekoch/gouniversal/module/picturex/pairlist"
	"github.com/dekoch/gouniversal/module/picturex/request"
	"github.com/dekoch/gouniversal/module/picturex/ui"
	"github.com/dekoch/gouniversal/module/picturex/upload"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	global.PairList.SetMaxPairs(global.Config.MaxPairs)
	pairlist.SetSourcePath(global.Config.RawFileRoot)
	pairlist.SetDestinationPath(global.Config.TempFileRoot)
	pairlist.SetStaticPath(global.Config.StaticFileRoot)

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
