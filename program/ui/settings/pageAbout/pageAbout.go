package pageAbout

import (
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Settings:About", page.Lang.Settings.About.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	page.Content += "written with Go"
}
