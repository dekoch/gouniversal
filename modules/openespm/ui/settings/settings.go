package settings

import (
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings/app"
	"gouniversal/modules/openespm/ui/settings/device"
	"gouniversal/shared/navigation"
	"net/http"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	app.RegisterPage(page, nav)
	device.RegisterPage(page, nav)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings" {
		nav.NavigatePath("App:openESPM:Settings:App")
	}

	if nav.IsNext("App") {

		app.Render(page, nav, r)

	} else if nav.IsNext("Device") {

		device.Render(page, nav, r)
	}
}
