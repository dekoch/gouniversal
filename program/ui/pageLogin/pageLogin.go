package pageLogin

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Account:Login", page.Lang.Menu.Account.Login)
	nav.Sitemap.Register("Account:Logout", page.Lang.Menu.Account.Logout)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type login struct {
		Lang lang.Login
	}
	var l login

	l.Lang = page.Lang.Login

	templ, err := template.ParseFiles(global.UiConfig.ProgramFileRoot + "login.html")
	if err != nil {
		fmt.Println(err)
	}
	page.Content += functions.TemplToString(templ, l)
}
