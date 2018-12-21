package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/fileshare/global"
	"github.com/dekoch/gouniversal/module/fileshare/typefileshare"
	"github.com/dekoch/gouniversal/module/fileshare/ui/pageedit"
	"github.com/dekoch/gouniversal/module/fileshare/ui/pagehome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	pagehome.LoadConfig()
	pageedit.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typefileshare.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagehome.RegisterPage(appPage, nav)
	pageedit.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typefileshare.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pagehome.Render(appPage, nav, r)

	} else if nav.IsNext("Edit") {

		pageedit.Render(appPage, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
