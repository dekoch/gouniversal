package pageGroup

import (
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/pageSettings/pageGroupEdit"
	"gouniversal/program/ui/pageSettings/pageGroupList"
	"gouniversal/program/ui/uiglobal"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:Group", page.Lang.Settings.Group.Title)
	pageGroupList.RegisterPage(page, nav)
	pageGroupEdit.RegisterPage(page, nav)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:Group" {
		nav.NavigatePath("Program:Settings:Group:List")
	}

	if nav.IsNext("List") {

		pageGroupList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageGroupEdit.Render(page, nav, r)
	}
}
