package settings

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/app"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/device"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typeoespm.Page, nav *navigation.Navigation) {

	app.RegisterPage(page, nav)
	device.RegisterPage(page, nav)
}

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings" {
		nav.NavigatePath("App:openESPM:Settings:App")
	}

	switch nav.GetNextPage() {
	case "App":
		app.Render(page, nav, r)

	case "Device":
		device.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
