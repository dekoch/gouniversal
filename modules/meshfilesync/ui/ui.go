package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/meshfilesync/global"
	"github.com/dekoch/gouniversal/modules/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/modules/meshfilesync/ui/pagesearch"
	"github.com/dekoch/gouniversal/modules/meshfilesync/ui/pagesettings"
	"github.com/dekoch/gouniversal/modules/meshfilesync/ui/pagetransfers"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesmfs.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagesearch.RegisterPage(appPage, nav)
	pagetransfers.RegisterPage(appPage, nav)
	pagesettings.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesmfs.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Search") {

		pagesearch.Render(appPage, nav, r)

	} else if nav.IsNext("Transfers") {

		pagetransfers.Render(appPage, nav, r)

	} else if nav.IsNext("Settings") {

		pagesettings.Render(appPage, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
