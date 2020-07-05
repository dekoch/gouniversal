package settings

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/schedule"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/uiplc"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	uiplc.RegisterPage(page, nav)
	schedule.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:S7Backup:Settings" {
		nav.NavigatePath("App:S7Backup:PLC")
	}

	switch nav.GetNextPage() {
	case "PLC":
		uiplc.Render(page, nav, r)

	case "Schedule":
		schedule.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
