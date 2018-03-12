package pageSettings

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui/pageSettings/pageAbout"
	"gouniversal/program/ui/pageSettings/pageGeneral"
	"gouniversal/program/ui/pageSettings/pageGroup"
	"gouniversal/program/ui/pageSettings/pageUser"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	//nav.Sitemap.Register("Program:Settings", page.Lang.Settings.Title)
	pageGeneral.RegisterPage(page, nav)
	pageUser.RegisterPage(page, nav)
	pageGroup.RegisterPage(page, nav)
	pageAbout.RegisterPage(page, nav)
}

func showTop(page *types.Page, nav *navigation.Navigation) {

	type sidebar struct {
		Lang  lang.Settings
		Title string
	}
	var sb sidebar

	sb.Lang = page.Lang.Settings
	sb.Title = nav.Sitemap.PageTitle(nav.Path)

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "/settings/settingstop.html")
	if err != nil {
		fmt.Println(err)
	}
	page.Content += functions.TemplToString(templ, sb)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	//showTop(page, nav)

	if nav.IsNext("General") {

		pageGeneral.Render(page, nav, r)

	} else if nav.IsNext("User") {

		pageUser.Render(page, nav, r)

	} else if nav.IsNext("Group") {

		pageGroup.Render(page, nav, r)

	} else if nav.IsNext("About") {

		pageAbout.Render(page, nav, r)
	}

	/*type sidebar struct {
		Title string
	}
	var sb sidebar

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "/settings/settingsbottom.html")
	if err != nil {
		fmt.Println(err)
	}
	page.Content += functions.TemplToString(templ, sb)*/
}
