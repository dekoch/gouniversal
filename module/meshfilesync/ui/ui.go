package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/meshfilesync/global"
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/module/meshfilesync/ui/pagesearch"
	"github.com/dekoch/gouniversal/module/meshfilesync/ui/pagesettings"
	"github.com/dekoch/gouniversal/module/meshfilesync/ui/pagetransfers"
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

	switch nav.GetNextPage() {
	case "Search":
		pagesearch.Render(appPage, nav, r)

	case "Transfers":
		pagetransfers.Render(appPage, nav, r)

	case "Settings":
		pagesettings.Render(appPage, nav, r)

	default:
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
