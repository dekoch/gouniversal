package ui

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/module/monmotion/ui/pageacquire"
	"github.com/dekoch/gouniversal/module/monmotion/ui/pagedevicelist"
	"github.com/dekoch/gouniversal/module/monmotion/ui/pagetrigger"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typemd.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	pagedevicelist.RegisterPage(appPage, nav)
	pageacquire.RegisterPage(appPage, nav)
	pagetrigger.RegisterPage(appPage, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typemd.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	switch nav.GetNextPage() {
	case "DeviceList":
		pagedevicelist.Render(appPage, nav, r)

	case "Acquire":
		pageacquire.Render(appPage, nav, r)

	case "Trigger":
		pagetrigger.Render(appPage, nav, r)

	default:
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}
