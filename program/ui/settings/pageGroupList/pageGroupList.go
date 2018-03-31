package pageGroupList

import (
	"fmt"
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

	nav.Sitemap.Register("Program:Settings:Group:List", page.Lang.Settings.Group.GroupList.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type grouplist struct {
		Lang      lang.SettingsGroupList
		GroupList template.HTML
	}
	var gl grouplist

	gl.Lang = page.Lang.Settings.Group.GroupList

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
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:Group:Edit$UUID=" + global.GroupConfig.File.Group[i].UUID + "\" title=\"" + gl.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	global.GroupConfig.Mut.Unlock()

	gl.GroupList = template.HTML(tbody)

	templ, err := template.ParseFiles(global.UiConfig.ProgramFileRoot + "settings/grouplist.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, gl)
}
