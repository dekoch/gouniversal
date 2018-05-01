package group

import (
	"gouniversal/program/ui/settings/group/pageGroupEdit"
	"gouniversal/program/ui/settings/group/pageGroupList"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Settings:Group", page.Lang.Settings.Group.Title)
	pageGroupList.RegisterPage(page, nav)
	pageGroupEdit.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings:Group" {
		nav.NavigatePath("Program:Settings:Group:List")
	}

	if nav.IsNext("List") {

		pageGroupList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageGroupEdit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
