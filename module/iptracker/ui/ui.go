package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/iptracker/global"
	"github.com/dekoch/gouniversal/module/iptracker/typeiptracker"
	"github.com/dekoch/gouniversal/module/iptracker/ui/pagehome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	pagehome.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typeiptracker.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagehome.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typeiptracker.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pagehome.Render(appPage, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
