package pageAbout

import (
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uiglobal"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:About", page.Lang.Settings.About.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	page.Content += "written with Go"
}
