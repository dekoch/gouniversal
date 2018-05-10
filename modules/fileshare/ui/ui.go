package ui

import (
	"gouniversal/modules/fileshare/global"
	"gouniversal/modules/fileshare/typesFileshare"
	"gouniversal/modules/fileshare/ui/home"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesFileshare.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	home.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesFileshare.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		home.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}

func LoadConfig() {

	home.LoadConfig()
}
