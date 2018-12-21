package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/logviewer/global"
	"github.com/dekoch/gouniversal/module/logviewer/typelogviewer"
	"github.com/dekoch/gouniversal/module/logviewer/ui/pagehome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typelogviewer.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagehome.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typelogviewer.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pagehome.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
