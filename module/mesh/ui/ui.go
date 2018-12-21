package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/module/mesh/ui/pagenetwork"
	"github.com/dekoch/gouniversal/module/mesh/ui/pageserver"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typemesh.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageserver.RegisterPage(appPage, nav)
	pagenetwork.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typemesh.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Server") {

		pageserver.Render(appPage, nav, r)
	} else if nav.IsNext("Network") {

		pagenetwork.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
