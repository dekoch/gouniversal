package station

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/module/gasprice/ui/settings/station/pagestationedit"
	"github.com/dekoch/gouniversal/module/gasprice/ui/settings/station/pagestationlist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	pagestationedit.RegisterPage(page, nav)
	pagestationlist.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:GasPrice:Settings:Station" {
		nav.NavigatePath("App:GasPrice:Station:List")
	}

	if nav.IsNext("Edit") {

		pagestationedit.Render(page, nav, r)

	} else if nav.IsNext("List") {

		pagestationlist.Render(page, nav, r)

	} else {
		nav.RedirectPath("404", true)
	}
}
