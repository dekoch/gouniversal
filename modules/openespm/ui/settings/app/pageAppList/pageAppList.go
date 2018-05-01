package pageAppList

import (
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

	type content struct {
		Lang    langOESPM.SettingsAppList
		AppList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.App.List

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
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:openESPM:Settings:App:Edit$UUID=" + a.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	globalOESPM.AppConfig.Mut.Unlock()

	c.AppList = template.HTML(tbody)

	p, err := functions.PageToString(globalOESPM.UiConfig.AppFileRoot+"settings/applist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
