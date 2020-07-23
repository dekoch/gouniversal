package pagebackup

import (
	"html"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/s7backup/dbconfig"
	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/s7"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:S7Backup:Backup:Backup", page.Lang.Backup.BackupTitle)
	nav.Sitemap.Register("", "App:S7Backup:Backup:Restore", page.Lang.Backup.RestoreTitle)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang            lang.Backup
		PLCName         template.HTML
		PLCUUID         template.HTML
		PLCAddress      template.HTML
		PLCRack         template.HTML
		PLCSlot         template.HTML
		PLCMaxBackupCnt template.HTML
		Comment         template.HTML
		CommentHidden   template.HTMLAttr
		DBList          template.HTML
		BTNPLC          template.HTML
	}
	var c Content

	c.Lang = page.Lang.Backup

	var (
		err error
		pc  plcconfig.PLCConfig
	)

	restoreView := false

	if nav.CurrentPath == "App:S7Backup:Backup:Restore" {
		restoreView = true
	}

	_, err = pc.Load(global.Config.GetFileRoot(), nav.Parameter("UUID"))
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	comment := html.EscapeString(r.FormValue("comment"))
	c.Comment = template.HTML(comment)

	// Form input
	backup := r.FormValue("backup")

	switch backup {
	case "plc":
		err = s7.BackupPLC(comment, &pc)
		if err == nil {
			alert.Message(alert.SUCCESS, page.Lang.Alert.Success, page.Lang.Backup.SavedToDatabase, " ", nav.User.UUID)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}

	default:
		for i := range pc.DB {

			if backup == "db"+pc.DB[i].UUID {

				var l []string
				l = append(l, pc.DB[i].UUID)

				err = s7.BackupDB(l, comment, &pc)
				if err == nil {
					alert.Message(alert.SUCCESS, page.Lang.Alert.Success, "DB"+strconv.Itoa(pc.DB[i].DBNo)+" "+page.Lang.Backup.SavedToDatabase, " ", nav.User.UUID)
				} else {
					alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				}
			}
		}
	}

	restore := r.FormValue("restore")

	switch restore {
	case "plc":
		err = s7.RestorePLC(&pc)
		if err == nil {
			alert.Message(alert.SUCCESS, page.Lang.Alert.Success, page.Lang.Backup.RestoredFromDatabase, " ", nav.User.UUID)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}

	default:
		for i := range pc.DB {

			var dcs []dbconfig.DBConfig
			dcs, err = pc.GetBackups(global.Config.GetFileRoot(), pc.DB[i].DBNo)
			if err != nil {
				continue
			}

			for i := range dcs {

				if restore == "db"+strconv.Itoa(dcs[i].ID) {

					var l []int
					l = append(l, dcs[i].ID)

					err = s7.RestoreDB(l, &pc)
					if err == nil {
						alert.Message(alert.SUCCESS, page.Lang.Alert.Success, "DB"+strconv.Itoa(dcs[i].DBNo)+" ("+dcs[i].Backup.Format("2006-01-02 15:04:05")+") "+page.Lang.Backup.RestoredFromDatabase, " ", nav.User.UUID)
					} else {
						alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
					}
				}
			}
		}
	}

	c.PLCName = template.HTML(pc.Name)
	c.PLCAddress = template.HTML(pc.Address)
	c.PLCRack = template.HTML(strconv.Itoa(pc.Rack))
	c.PLCSlot = template.HTML(strconv.Itoa(pc.Slot))
	c.PLCMaxBackupCnt = template.HTML(strconv.Itoa(pc.MaxBackupCnt))
	c.PLCUUID = template.HTML(pc.UUID)

	if restoreView {
		c.CommentHidden = template.HTMLAttr(" hidden")
		c.BTNPLC = template.HTML("<button class=\"btn btn-default fa fa-microchip\" type=\"submit\" name=\"restore\" value=\"plc\" title=\"" + c.Lang.RestorePLC + "\"></button>")
	} else {
		c.CommentHidden = template.HTMLAttr("")
		c.BTNPLC = template.HTML("<button class=\"btn btn-default fa fa-database\" type=\"submit\" name=\"backup\" value=\"plc\" title=\"" + c.Lang.BackupPLC + "\"></button>")
	}

	dbList := ""

	for _, db := range pc.DB {

		if db.DBNo <= 0 ||
			db.DBByteLength <= 0 {

			continue
		}

		dbList += "<tr>"
		dbList += "<td>" + db.Name + "</td>"
		dbList += "<td>" + strconv.Itoa(db.DBNo) + "</td>"
		dbList += "<td>" + strconv.Itoa(db.DBByteLength) + "</td>"

		if restoreView {

			var dcs []dbconfig.DBConfig
			dcs, err = pc.GetBackups(global.Config.GetFileRoot(), db.DBNo)
			if err != nil {
				continue
			}

			dbList += "<td>"

			for i := range dcs {
				dbList += "<button class=\"btn btn-default fa fa-microchip\" type=\"submit\" name=\"restore\" value=\"db" + strconv.Itoa(dcs[i].ID) + "\" title=\"" + c.Lang.RestoreDB + "\">" + dcs[i].Backup.Format("2006-01-02 15:04:05") + "<br>" + dcs[i].Comment + "</button><br>"
			}

			dbList += "</td>"
		} else {
			dbList += "<td><button class=\"btn btn-default fa fa-database\" type=\"submit\" name=\"backup\" value=\"db" + db.UUID + "\" title=\"" + c.Lang.BackupDB + "\"></button></td>"
		}
		dbList += "</tr>"
	}

	c.DBList = template.HTML(dbList)

	p, err := functions.PageToString(global.Config.UIFileRoot+"backupplc.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
