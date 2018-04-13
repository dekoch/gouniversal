package pageAppList

import (
	"fmt"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"net/http"
	"strconv"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:App:List", page.Lang.Settings.App.List.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	type appList struct {
		Lang    langOESPM.SettingsAppList
		AppList template.HTML
	}
	var al appList

	al.Lang = page.Lang.Settings.App.List

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	globalOESPM.AppConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.AppConfig.File.Apps); i++ {

		a := globalOESPM.AppConfig.File.Apps[i]

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + a.Name + "</td>"
		tbody += "<td>" + a.App + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:openESPM:Settings:App:Edit$UUID=" + a.UUID + "\" title=\"" + al.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	globalOESPM.AppConfig.Mut.Unlock()

	al.AppList = template.HTML(tbody)

	templ, err := template.ParseFiles(globalOESPM.UiConfig.AppFileRoot + "settings/applist.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, al)
}
