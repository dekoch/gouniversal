package ui

import (
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("App:openESPM", "openESPM")
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	page.Content = "hello from openESPM"
}
