package user

import (
	"net/http"

	"github.com/dekoch/gouniversal/program/ui/settings/user/pageuseredit"
	"github.com/dekoch/gouniversal/program/ui/settings/user/pageuserlist"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Settings:User", page.Lang.Settings.User.Title)
	pageuserlist.RegisterPage(page, nav)
	pageuseredit.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:User" {
		nav.NavigatePath("Program:Settings:User:List")
	}

	switch nav.GetNextPage() {
	case "List":
		pageuserlist.Render(page, nav, r)

	case "Edit":
		pageuseredit.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
