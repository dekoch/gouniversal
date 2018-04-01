package settings

import (
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings/device"
	"gouniversal/shared/navigation"
	"net/http"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	device.RegisterPage(page, nav)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:Program:openESPM:Settings" {
		nav.NavigatePath("App:Program:openESPM:Settings:Device")
	}

	if nav.IsNext("Device") {

		device.Render(page, nav, r)
	}
}
