package iptracker

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/iptracker/global"
	"github.com/dekoch/gouniversal/module/iptracker/lang"
	"github.com/dekoch/gouniversal/module/iptracker/tracker"
	"github.com/dekoch/gouniversal/module/iptracker/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/iptracker/", en, "en")

	global.Config.LoadConfig()

	tracker.LoadConfig()

	ui.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	global.Config.SaveConfig()
}
