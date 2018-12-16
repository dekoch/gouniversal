package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/typesMesh"
	"github.com/dekoch/gouniversal/module/mesh/ui/pageNetwork"
	"github.com/dekoch/gouniversal/module/mesh/ui/pageServer"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesMesh.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageServer.RegisterPage(appPage, nav)
	pageNetwork.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesMesh.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Server") {

		pageServer.Render(appPage, nav, r)
	} else if nav.IsNext("Network") {

		pageNetwork.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
