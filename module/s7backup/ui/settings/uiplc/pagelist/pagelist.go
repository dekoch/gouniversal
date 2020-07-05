package pagelist

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.PLCList.Menu, "App:S7Backup:Settings:PLC:List", page.Lang.PLCList.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang    lang.PLCList
		PLCList template.HTML
	}
	var c Content

	c.Lang = page.Lang.PLCList

	files, err := plcconfig.GetPLCList(global.Config.GetFileRoot())
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	tbody := ""

	for _, file := range files {

		var pc plcconfig.PLCConfig
		_, err := pc.Load(global.Config.GetFileRoot(), file)
		if err != nil {
			continue
		}

		tbody += "<tr>"
		tbody += "<td>" + pc.Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:S7Backup:Settings:PLC:Edit$UUID=" + pc.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.PLCList = template.HTML(tbody)

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/plclist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
