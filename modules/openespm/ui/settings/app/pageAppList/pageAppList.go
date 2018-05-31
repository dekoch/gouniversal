package pageAppList

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/modules/openespm/global"
	"github.com/dekoch/gouniversal/modules/openespm/lang"
	"github.com/dekoch/gouniversal/modules/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:App:List", page.Lang.Settings.App.List.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang    lang.SettingsAppList
		AppList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.App.List

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	apps := global.AppConfig.List()

	for i := 0; i < len(apps); i++ {

		a := apps[i]

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + a.Name + "</td>"
		tbody += "<td>" + a.App + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:openESPM:Settings:App:Edit$UUID=" + a.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.AppList = template.HTML(tbody)

	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"settings/applist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
