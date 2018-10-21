package pageUserList

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:User:List", page.Lang.Settings.User.UserList.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang     lang.SettingsUserList
		UserList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.User.UserList

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	user := global.UserConfig.List()

	for i := 0; i < len(user); i++ {

		u := user[i]

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + u.LoginName + "</td>"
		tbody += "<td>" + u.Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:User:Edit$UUID=" + u.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.UserList = template.HTML(tbody)

	p, err := functions.PageToString(global.UiConfig.ProgramFileRoot+"settings/userlist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
