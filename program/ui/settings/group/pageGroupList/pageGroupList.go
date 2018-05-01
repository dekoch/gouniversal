package pageGroupList

import (
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
	"strconv"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:Group:List", page.Lang.Settings.Group.GroupList.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang      lang.SettingsGroupList
		GroupList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Group.GroupList

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	global.GroupConfig.Mut.Lock()
	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + global.GroupConfig.File.Group[i].Name + "</td>"
		tbody += "<td>" + global.GroupConfig.File.Group[i].Comment + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:Group:Edit$UUID=" + global.GroupConfig.File.Group[i].UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	global.GroupConfig.Mut.Unlock()

	c.GroupList = template.HTML(tbody)

	p, err := functions.PageToString(global.UiConfig.File.ProgramFileRoot+"settings/grouplist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
