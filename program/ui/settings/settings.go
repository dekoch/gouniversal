package settings

import (
	"gouniversal/program/ui/settings/pageAbout"
	"gouniversal/program/ui/settings/pageGeneral"
	"gouniversal/program/ui/settings/pageGroup"
	"gouniversal/program/ui/settings/pageUser"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	pageGeneral.RegisterPage(page, nav)
	pageUser.RegisterPage(page, nav)
	pageGroup.RegisterPage(page, nav)
	pageAbout.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings" {
		nav.NavigatePath("Program:Settings:General")
	}

	if nav.IsNext("General") {

		pageGeneral.Render(page, nav, r)

	} else if nav.IsNext("User") {

		pageUser.Render(page, nav, r)

	} else if nav.IsNext("Group") {

		pageGroup.Render(page, nav, r)

	} else if nav.IsNext("About") {

		pageAbout.Render(page, nav, r)
	}
}
