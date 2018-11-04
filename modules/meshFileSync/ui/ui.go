package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/typesMFS"
	"github.com/dekoch/gouniversal/modules/meshFileSync/ui/pageSearch"
	"github.com/dekoch/gouniversal/modules/meshFileSync/ui/pageSettings"
	"github.com/dekoch/gouniversal/modules/meshFileSync/ui/pageTransfers"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesMFS.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageSearch.RegisterPage(appPage, nav)
	pageTransfers.RegisterPage(appPage, nav)
	pageSettings.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesMFS.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Search") {

		pageSearch.Render(appPage, nav, r)

	} else if nav.IsNext("Transfers") {

		pageTransfers.Render(appPage, nav, r)

	} else if nav.IsNext("Settings") {

		pageSettings.Render(appPage, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
