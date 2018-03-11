package pageUserList

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"html/template"
	"net/http"
	"strconv"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:User:List", page.Lang.Settings.User.UserList.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	type userlist struct {
		Lang     lang.SettingsUserList
		UserList template.HTML
	}
	var ul userlist

	ul.Lang = page.Lang.Settings.User.UserList

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	global.UserConfig.Mut.Lock()
	for i := 0; i < len(global.UserConfig.File.User); i++ {

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + global.UserConfig.File.User[i].LoginName + "</td>"
		tbody += "<td>" + global.UserConfig.File.User[i].Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:User:Edit$UUID=" + global.UserConfig.File.User[i].UUID + "\" title=\"" + ul.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	global.UserConfig.Mut.Unlock()

	ul.UserList = template.HTML(tbody)

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/settings/userlist.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += uifunc.TemplToString(templ, ul)
}
