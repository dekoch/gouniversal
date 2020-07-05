package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/module/s7backup/ui/backup"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typemd.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	backup.RegisterPage(appPage, nav)
	settings.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typemd.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	switch nav.GetNextPage() {
	case "Backup":
		backup.Render(appPage, nav, r)

	case "Settings":
		settings.Render(appPage, nav, r)

	default:
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
