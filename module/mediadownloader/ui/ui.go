package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/mediadownloader/global"
	"github.com/dekoch/gouniversal/module/mediadownloader/typesMD"
	"github.com/dekoch/gouniversal/module/mediadownloader/ui/pageHome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesMD.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageHome.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesMD.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pageHome.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
