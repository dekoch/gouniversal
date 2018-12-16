package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/console/global"
	"github.com/dekoch/gouniversal/module/console/typesConsole"
	"github.com/dekoch/gouniversal/module/console/ui/pageHome"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	pageHome.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesConsole.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pageHome.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesConsole.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Home") {

		pageHome.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
