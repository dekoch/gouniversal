package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/backup/global"
	"github.com/dekoch/gouniversal/module/backup/typemd"
	"github.com/dekoch/gouniversal/module/backup/ui/pagehome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typemd.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagehome.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typemd.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	switch nav.GetNextPage() {
	case "Home":
		pagehome.Render(appPage, nav, r)

	default:
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
