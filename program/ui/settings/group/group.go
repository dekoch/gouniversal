package group

import (
	"net/http"

	"github.com/dekoch/gouniversal/program/ui/settings/group/pagegroupedit"
	"github.com/dekoch/gouniversal/program/ui/settings/group/pagegrouplist"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Settings:Group", page.Lang.Settings.Group.Title)
	pagegrouplist.RegisterPage(page, nav)
	pagegroupedit.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:Group" {
		nav.NavigatePath("Program:Settings:Group:List")
	}

	if nav.IsNext("List") {

		pagegrouplist.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pagegroupedit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
