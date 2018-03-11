package pageGeneral

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/programConfig"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"html/template"
	"net/http"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:General", page.Lang.Settings.GeneralEdit.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	var strButton string

	strButton = r.FormValue("edit")

	fmt.Println(strButton)

	if strButton == "apply" {

		programConfig.SaveConfig(global.ProgramConfig)
	}

	type general struct {
		Lang lang.SettingsGeneral
	}
	var g general

	g.Lang = page.Lang.Settings.GeneralEdit

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/settings/general.html")
	if err != nil {
		fmt.Println(err)
	}
	page.Content += uifunc.TemplToString(templ, g)

}
