package pagelist

import (
	"html/template"
	"net/http"
	"os"

	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/moduleconfig/scheduleconfig"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.ScheduleList.Menu, "App:S7Backup:Settings:Schedule:List", page.Lang.ScheduleList.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang                lang.ScheduleList
		ScheduleList        template.HTML
		CmbPLC              template.HTML
		AddScheduleDisabled template.HTMLAttr
	}
	var c Content

	c.Lang = page.Lang.ScheduleList

	switch r.FormValue("edit") {
	case "addschedule":
		uid, err := addSchedule(r)
		if err == nil {
			nav.RedirectPath("App:S7Backup:Settings:Schedule:Edit$UUID="+uid, false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	tbody := ""

	for _, schedule := range global.Config.Schedule.GetList() {

		if _, err := os.Stat(global.Config.GetFileRoot() + schedule.PLC + ".sqlite3"); os.IsNotExist(err) {
			continue
		}

		var pc plcconfig.PLCConfig
		_, err := pc.Load(global.Config.GetFileRoot(), schedule.PLC)
		if err != nil {
			continue
		}

		tbody += "<tr>"
		tbody += "<td>" + schedule.Name + "</td>"
		tbody += "<td>" + pc.Name + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:S7Backup:Settings:Schedule:Edit$UUID=" + schedule.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.ScheduleList = template.HTML(tbody)

	cmb, err := cmbPLC()
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	c.CmbPLC = template.HTML("<select name=\"cmbplc\">" + cmb + "</select>")

	if cmb == "" {
		c.AddScheduleDisabled = template.HTMLAttr(" disabled")
	} else {
		c.AddScheduleDisabled = template.HTMLAttr("")
	}

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/schedulelist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func addSchedule(r *http.Request) (string, error) {

	var (
		err error
		ret string
	)

	func() {

		var (
			cmbPLC string
			files  []string
		)

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				cmbPLC, err = functions.CheckFormInput("cmbplc", r)
				if functions.IsEmpty(cmbPLC) {
					return
				}

			case 1:
				files, err = plcconfig.GetPLCList(global.Config.GetFileRoot())

			case 2:
				for i := range files {

					if files[i] == cmbPLC {

						bs := scheduleconfig.NewBackupSchedule()
						bs.PLC = cmbPLC

						ret = bs.UUID
						err = global.Config.Schedule.Add(bs)
					}
				}

			case 3:
				err = global.Config.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func cmbPLC() (string, error) {

	files, err := plcconfig.GetPLCList(global.Config.GetFileRoot())
	if err != nil {
		return "", err
	}

	var ret string

	for _, file := range files {

		var pc plcconfig.PLCConfig
		_, err = pc.Load(global.Config.GetFileRoot(), file)
		if err != nil {
			continue
		}

		ret += "<option value=\"" + file + "\">" + pc.Name + "</option>"
	}

	return ret, nil
}
