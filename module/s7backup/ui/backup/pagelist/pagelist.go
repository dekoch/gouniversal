package pagelist

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/program/usermanagement"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.BackupList.Menu, "App:S7Backup:Backup:List", page.Lang.BackupList.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang    lang.BackupList
		PLCList template.HTML
	}
	var c Content

	c.Lang = page.Lang.BackupList

	var (
		backupDisabled  string
		restoreDisabled string
	)

	if usermanagement.IsPageAllowed("App:S7Backup:Backup:Backup", nav.User) == false {
		backupDisabled = " disabled"
	}

	if usermanagement.IsPageAllowed("App:S7Backup:Backup:Restore", nav.User) == false {
		restoreDisabled = " disabled"
	}

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

		if len(pc.DB) == 0 {
			continue
		}

		tbody += "<tr>"
		tbody += "<td>" + pc.Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-database\" type=\"submit\" name=\"navigation\" value=\"App:S7Backup:Backup:Backup$UUID=" + pc.UUID + "\" title=\"" + c.Lang.Backup + "\"" + backupDisabled + "></button></td>"
		tbody += "<td><button class=\"btn btn-default fa fa-microchip\" type=\"submit\" name=\"navigation\" value=\"App:S7Backup:Backup:Restore$UUID=" + pc.UUID + "\" title=\"" + c.Lang.Restore + "\"" + restoreDisabled + "></button></td>"
		tbody += "</tr>"
	}

	c.PLCList = template.HTML(tbody)

	p, err := functions.PageToString(global.Config.UIFileRoot+"backuplist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
