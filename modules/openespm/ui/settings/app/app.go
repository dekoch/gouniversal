package app

import (
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings/app/pageAppEdit"
	"gouniversal/modules/openespm/ui/settings/app/pageAppList"
	"gouniversal/shared/navigation"
	"net/http"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("openESPM", "App:openESPM:Settings:App", page.Lang.Settings.App.Title)
	pageAppList.RegisterPage(page, nav)
	pageAppEdit.RegisterPage(page, nav)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings:App" {
		nav.NavigatePath("App:openESPM:Settings:App:List")
	}

	if nav.IsNext("List") {

		pageAppList.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageAppEdit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
