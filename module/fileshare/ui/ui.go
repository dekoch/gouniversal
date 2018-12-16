package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/fileshare/global"
	"github.com/dekoch/gouniversal/module/fileshare/typesFileshare"
	"github.com/dekoch/gouniversal/module/fileshare/ui/pageEdit"
	"github.com/dekoch/gouniversal/module/fileshare/ui/pageHome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	pageHome.LoadConfig()
	pageEdit.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesFileshare.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageHome.RegisterPage(appPage, nav)
	pageEdit.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesFileshare.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pageHome.Render(appPage, nav, r)

	} else if nav.IsNext("Edit") {

		pageEdit.Render(appPage, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
