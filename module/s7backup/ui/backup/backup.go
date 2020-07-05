package backup

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/module/s7backup/ui/backup/pagebackup"
	"github.com/dekoch/gouniversal/module/s7backup/ui/backup/pagelist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	pagelist.RegisterPage(page, nav)
	pagebackup.RegisterPage(page, nav)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:S7Backup:Backup" {
		nav.NavigatePath("App:S7Backup:Backup:List")
	}

	switch nav.GetNextPage() {
	case "List":
		pagelist.Render(page, nav, r)

	case "Backup":
		pagebackup.Render(page, nav, r)

	case "Restore":
		pagebackup.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
