package pageHome

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Home", page.Lang.Home.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	templ, err := template.ParseFiles(global.UiConfig.ProgramFileRoot + "home.html")
	if err != nil {
		fmt.Println(err)
	}

	items := struct {
		Temp string
	}{
		Temp: "",
	}
	page.Content += functions.TemplToString(templ, items)
}
