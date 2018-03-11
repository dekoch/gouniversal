package pageLogin

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"html/template"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Login", page.Lang.Login.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	page.Title = page.Lang.Login.Title

	type login struct {
		Lang lang.Login
	}
	var l login

	l.Lang = page.Lang.Login

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/login.html")
	if err != nil {
		fmt.Println(err)
	}
	page.Content += uifunc.TemplToString(templ, l)
}
