package pageuserlist

import (
	"html/template"
	"net/http"
	"sort"
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

	tbody := ""

	user := global.UserConfig.List()

	sort.Slice(user, func(i, j int) bool { return user[i].LoginName < user[j].LoginName })

	for i, u := range user {

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(i+1) + "</th>"
		tbody += "<td>" + u.LoginName + "</td>"
		tbody += "<td>" + u.Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:User:Edit$UUID=" + u.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.UserList = template.HTML(tbody)

	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"settings/userlist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
