package app

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/app/pageappedit"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings/app/pageapplist"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typeoespm.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("openESPM", "App:openESPM:Settings:App", page.Lang.Settings.App.Title)
	pageapplist.RegisterPage(page, nav)
	pageappedit.RegisterPage(page, nav)
}

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "App:openESPM:Settings:App" {
		nav.NavigatePath("App:openESPM:Settings:App:List")
	}

	if nav.IsNext("List") {

		pageapplist.Render(page, nav, r)

	} else if nav.IsNext("Edit") {

		pageappedit.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
