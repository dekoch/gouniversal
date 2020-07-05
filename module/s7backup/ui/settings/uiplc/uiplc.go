package uiplc

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/uiplc/pageedit"
	"github.com/dekoch/gouniversal/module/s7backup/ui/settings/uiplc/pagelist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	pageedit.RegisterPage(page, nav)
	pagelist.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:S7Backup:Settings:PLC" {
		nav.NavigatePath("App:S7Backup:PLC:List")
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
