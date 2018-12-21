package device

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/device/pagedeviceedit"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/device/pagedevicelist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typeoespm.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("openESPM", "App:openESPM:Settings:Device", page.Lang.Settings.Device.Title)
	pagedevicelist.RegisterPage(page, nav)
	pagedeviceedit.RegisterPage(page, nav)
}

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings:Device" {
		nav.NavigatePath("App:openESPM:Settings:Device:List")
	}

	if nav.IsNext("List") {

		pagedevicelist.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pagedeviceedit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
