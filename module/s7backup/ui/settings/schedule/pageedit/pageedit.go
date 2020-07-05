package pageedit

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

	nav.Sitemap.Register("", "App:S7Backup:Settings:Schedule:Edit", page.Lang.ScheduleEdit.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang         lang.ScheduleEdit
		ScheduleName template.HTML
		PLCName      template.HTML
		CheckedSu    template.HTMLAttr
		CheckedMo    template.HTMLAttr
		CheckedTu    template.HTMLAttr
		CheckedWe    template.HTMLAttr
		CheckedTh    template.HTMLAttr
		CheckedFr    template.HTMLAttr
		CheckedSa    template.HTMLAttr
		ScheduleUUID template.HTML
		DBList       template.HTML
	}
	var c Content

	c.Lang = page.Lang.ScheduleEdit

	// Form input
	bs, err := global.Config.Schedule.Get(nav.Parameter("UUID"))
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	var pc plcconfig.PLCConfig
	_, err = pc.Load(global.Config.GetFileRoot(), bs.PLC)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	switch r.FormValue("edit") {
	case "apply":
		err = editSchedule(r, &bs)
		if err == nil {
			nav.RedirectPath("App:S7Backup:Settings:Schedule:List", false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)

	case "delete":
		err = deleteSchedule(bs.UUID)
		if err == nil {
			nav.RedirectPath("App:S7Backup:Settings:Schedule:List", false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	c.ScheduleName = template.HTML(bs.Name)
	c.PLCName = template.HTML(pc.Name)
	c.CheckedSu = isDayChecked(0, &bs)
	c.CheckedMo = isDayChecked(1, &bs)
	c.CheckedTu = isDayChecked(2, &bs)
	c.CheckedWe = isDayChecked(3, &bs)
	c.CheckedTh = isDayChecked(4, &bs)
	c.CheckedFr = isDayChecked(5, &bs)
	c.CheckedSa = isDayChecked(6, &bs)
	c.ScheduleUUID = template.HTML(bs.UUID)

	dbList := ""

	for _, db := range pc.DB {

		if isDBChecked(db.UUID, &bs) {
			dbList += "<tr class=\"table-success\">"
		} else {
			dbList += "<tr>"
		}

		dbList += "<td><input type=\"checkbox\" name=\"selecteddb\" value=\"" + db.UUID + "\""

		if isDBChecked(db.UUID, &bs) {
			dbList += " checked"
		}

		dbList += "></td>"

		dbList += "<td>" + db.Name + "</td>"
		dbList += "<td>" + strconv.Itoa(db.DBNo) + "</td>"
		dbList += "</tr>"
	}

	c.DBList = template.HTML(dbList)

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/scheduleedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func editSchedule(r *http.Request, bs *scheduleconfig.BackupSchedule) error {

	var err error

	func() {

		var (
			name string
		)

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				name, err = functions.CheckFormInput("schedulename", r)

			case 1:
				if functions.IsEmpty(name) {

					err = errors.New("bad input")
				}

			case 2:
				bs.DB = r.Form["selecteddb"]

			case 3:
				days := r.Form["selectedday"]

				for i := range bs.Day {

					bs.Day[i] = false

					for ii := range days {

						var no int
						no, err = strconv.Atoi(days[ii])
						if err != nil {
							return
						}

						if i == no {
							bs.Day[i] = true
						}
					}
				}

			case 4:
				bs.Name = name
				bs.Backup = time.Now().AddDate(-999, 0, 0)
				bs.SetChecked(time.Now().AddDate(-999, 0, 0))
				err = global.Config.Schedule.Add(*bs)

			case 5:
				err = global.Config.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteSchedule(uid string) error {

	var err error

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				err = global.Config.Schedule.Delete(uid)

			case 1:
				err = global.Config.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func isDayChecked(day int, bs *scheduleconfig.BackupSchedule) template.HTMLAttr {

	if day > len(bs.Day) {
		return template.HTMLAttr("")
	}

	if bs.Day[day] {
		return template.HTMLAttr(" checked")
	}

	return template.HTMLAttr("")
}

func isDBChecked(uid string, bs *scheduleconfig.BackupSchedule) bool {

	for i := range bs.DB {

		if bs.DB[i] == uid {
			return true
		}
	}

	return false
}
