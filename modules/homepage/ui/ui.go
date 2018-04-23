package ui

import (
	"fmt"
	"gouniversal/modules/homepage/globalHomepage"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Home", "App:homepage:Home", "Home")
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Index template.HTML
	}
	var c Content

	templ, err := template.ParseFiles(globalHomepage.ContentFolder + "index.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, c)
}
