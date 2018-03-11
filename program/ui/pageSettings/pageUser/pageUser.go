package pageUser

import (
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/pageSettings/pageUserEdit"
	"gouniversal/program/ui/pageSettings/pageUserList"
	"gouniversal/program/ui/uiglobal"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:User", page.Lang.Settings.User.Title)
	pageUserList.RegisterPage(page, nav)
	pageUserEdit.RegisterPage(page, nav)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:User" {
		nav.NavigatePath("Program:Settings:User:List")
	}

	if nav.IsNext("List") {

		pageUserList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageUserEdit.Render(page, nav, r)
	}
}
