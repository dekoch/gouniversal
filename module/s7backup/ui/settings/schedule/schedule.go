package schedule

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/schedule/pageedit"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/schedule/pagelist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	pagelist.RegisterPage(page, nav)
	pageedit.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:S7Backup:Settings:Schedule" {
		nav.NavigatePath("App:S7Backup:Settings:Schedule:List")
	}

	switch nav.GetNextPage() {
	case "List":
		pagelist.Render(page, nav, r)

	case "Edit":
		pageedit.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
