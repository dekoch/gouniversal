package pageHome

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"html/template"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Home", page.Lang.Home.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	page.Title = page.Lang.Home.Title

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/home.html")
	if err != nil {
		fmt.Println(err)
	}

	items := struct {
		Temp string
	}{
		Temp: "",
	}
	page.Content += uifunc.TemplToString(templ, items)
}
