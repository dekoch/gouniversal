package device

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/device/pageDeviceEdit"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/device/pageDeviceList"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("openESPM", "App:openESPM:Settings:Device", page.Lang.Settings.Device.Title)
	pageDeviceList.RegisterPage(page, nav)
	pageDeviceEdit.RegisterPage(page, nav)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings:Device" {
		nav.NavigatePath("App:openESPM:Settings:Device:List")
	}

	if nav.IsNext("List") {

		pageDeviceList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageDeviceEdit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
