package pagegrouplist

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

	nav.Sitemap.Register("", "Program:Settings:Group:List", page.Lang.Settings.Group.GroupList.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang      lang.SettingsGroupList
		GroupList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Group.GroupList

	tbody := ""

	groups := global.GroupConfig.List()

	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })

	for i, g := range groups {

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(i+1) + "</th>"
		tbody += "<td>" + g.Name + "</td>"
		tbody += "<td>" + g.Comment + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"Program:Settings:Group:Edit$UUID=" + g.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.GroupList = template.HTML(tbody)

	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"settings/grouplist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
