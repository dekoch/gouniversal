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

	if nav.IsNext("List") {

		pageuserlist.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageuseredit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
