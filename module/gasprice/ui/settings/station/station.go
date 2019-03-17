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

	switch nav.GetNextPage() {
	case "Edit":
		pagestationedit.Render(page, nav, r)

	case "List":
		pagestationlist.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
