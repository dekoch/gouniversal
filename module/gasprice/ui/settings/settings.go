package settings

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/module/gasprice/ui/settings/station"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	station.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:GasPrice:Settings" {
		nav.NavigatePath("App:GasPrice:Station")
	}

	if nav.IsNext("Station") {

		station.Render(page, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}
}
