package device

import (
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings/device/pageDeviceEdit"
	"gouniversal/modules/openespm/ui/settings/device/pageDeviceList"
	"gouniversal/shared/navigation"
	"net/http"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("App:Program:openESPM:Settings:Device", page.Lang.Settings.Device.Title)
	pageDeviceList.RegisterPage(page, nav)
	pageDeviceEdit.RegisterPage(page, nav)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:Program:openESPM:Settings:Device" {
		nav.NavigatePath("App:Program:openESPM:Settings:Device:List")
	}

	if nav.IsNext("List") {

		pageDeviceList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageDeviceEdit.Render(page, nav, r)
	}
}
