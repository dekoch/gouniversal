package user

import (
	"gouniversal/program/ui/settings/user/pageUserEdit"
	"gouniversal/program/ui/settings/user/pageUserList"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Settings:User", page.Lang.Settings.User.Title)
	pageUserList.RegisterPage(page, nav)
	pageUserEdit.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:User" {
		nav.NavigatePath("Program:Settings:User:List")
	}

	if nav.IsNext("List") {

		pageUserList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageUserEdit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
