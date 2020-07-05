package pageedit

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/module/s7backup/dbconfig"
	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/lang"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:S7Backup:Settings:PLC:Edit", page.Lang.PLCEdit.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang            lang.PLCEdit
		PLCName         template.HTML
		PLCAddress      template.HTML
		PLCRack         template.HTML
		PLCSlot         template.HTML
		PLCMaxBackupCnt template.HTML
		PLCUUID         template.HTML
		DBList          template.HTML
	}
	var c Content

	c.Lang = page.Lang.PLCEdit

	var (
		err error
		pc  plcconfig.PLCConfig
	)

	// Form input
	id := nav.Parameter("UUID")

	if id == "new" {

		id, err = newPLC()
		if err == nil {
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)

	} else {
		_, err = pc.Load(global.Config.GetFileRoot(), id)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	switch r.FormValue("edit") {
	case "apply":
		err = editPLC(r, &pc)
		if err == nil {
			nav.RedirectPath("App:S7Backup:Settings:PLC:List", false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)

	case "delete":
		err = deletePLC(&pc)
		if err == nil {
			nav.RedirectPath("App:S7Backup:Settings:PLC:List", false)
			return
		}

		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)

	case "adddb":
		err = addDB(r, &pc)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	uid, err := functions.CheckFormInput("deletedb", r)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	if uid != "" {

		err = deleteDB(uid, r, &pc)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	c.PLCName = template.HTML(pc.Name)
	c.PLCAddress = template.HTML(pc.Address)
	c.PLCRack = template.HTML(strconv.Itoa(pc.Rack))
	c.PLCSlot = template.HTML(strconv.Itoa(pc.Slot))
	c.PLCMaxBackupCnt = template.HTML(strconv.Itoa(pc.MaxBackupCnt))
	c.PLCUUID = template.HTML(pc.UUID)

	dbList := ""

	for _, db := range pc.DB {

		dbList += "<tr>"
		dbList += "<td><input type=\"text\" class=\"form-control\" name=\"dbname" + db.UUID + "\" value=\"" + db.Name + "\"></td>"
		dbList += "<td><input type=\"text\" class=\"form-control\" name=\"dbno" + db.UUID + "\" value=\"" + strconv.Itoa(db.DBNo) + "\"></td>"
		dbList += "<td><input type=\"text\" class=\"form-control\" name=\"dbbytelength" + db.UUID + "\" value=\"" + strconv.Itoa(db.DBByteLength) + "\"></td>"
		dbList += "<td><button class=\"btn btn-default fa fa-trash\" type=\"submit\" name=\"deletedb\" value=\"" + db.UUID + "\" title=\"" + c.Lang.Delete + "\"></button></td>"
		dbList += "</tr>"
	}

	c.DBList = template.HTML(dbList)

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/plcedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newPLC() (string, error) {

	var (
		err error
		pc  plcconfig.PLCConfig
	)

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				pc, err = plcconfig.NewPLC()

			case 1:
				err = pc.Save(global.Config.GetFileRoot())
			}

			if err != nil {
				return
			}
		}
	}()

	return pc.UUID, err
}

func editPLC(r *http.Request, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			name            string
			address         string
			rack            string
			intRack         int
			slot            string
			intSlot         int
			maxBackupCnt    string
			intMaxBackupCnt int
		)

		for i := 0; i <= 11; i++ {

			switch i {
			case 0:
				name, err = functions.CheckFormInput("plcname", r)

			case 1:
				address, err = functions.CheckFormInput("plcaddress", r)

			case 2:
				rack, err = functions.CheckFormInput("plcrack", r)

			case 3:
				slot, err = functions.CheckFormInput("plcslot", r)

			case 4:
				maxBackupCnt, err = functions.CheckFormInput("plcmaxbackupcnt", r)

			case 5:
				if functions.IsEmpty(name) ||
					functions.IsEmpty(address) ||
					functions.IsEmpty(rack) ||
					functions.IsEmpty(slot) ||
					functions.IsEmpty(maxBackupCnt) {

					err = errors.New("bad input")
				}

			case 6:
				intRack, err = strconv.Atoi(rack)

			case 7:
				intSlot, err = strconv.Atoi(slot)

			case 8:
				intMaxBackupCnt, err = strconv.Atoi(maxBackupCnt)

			case 9:
				if intRack < 0 ||
					intSlot < 0 ||
					intMaxBackupCnt < 0 {
					err = errors.New("bad input")
				}

			case 10:
				for i := range pc.DB {

					err = editDB(&pc.DB[i], r, pc)
					if err != nil {
						return
					}
				}

			case 11:
				pc.Name = name
				pc.Address = address
				pc.Rack = intRack
				pc.Slot = intSlot
				pc.MaxBackupCnt = intMaxBackupCnt
				err = pc.Save(global.Config.GetFileRoot())
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deletePLC(pc *plcconfig.PLCConfig) error {

	return pc.Delete(global.Config.GetFileRoot())
}

func addDB(r *http.Request, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var n dbconfig.DBConfig

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				err = editPLC(r, pc)

			case 1:
				n, err = dbconfig.NewDB()

			case 2:
				err = pc.AddDB(n)

			case 3:
				err = pc.Save(global.Config.GetFileRoot())
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteDB(uid string, r *http.Request, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				err = editPLC(r, pc)

			case 1:
				err = pc.DeleteDB(uid)

			case 2:
				err = pc.Save(global.Config.GetFileRoot())
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func editDB(dc *dbconfig.DBConfig, r *http.Request, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			name    string
			strDBNo string
			//intDBNo int
			strDBByteLength string
			//intDBByteLength int
		)

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				name, err = functions.CheckFormInput("dbname"+dc.UUID, r)

			case 1:
				strDBNo, err = functions.CheckFormInput("dbno"+dc.UUID, r)

			case 2:
				strDBByteLength, err = functions.CheckFormInput("dbbytelength"+dc.UUID, r)

			case 3:
				if functions.IsEmpty(name) {
					err = errors.New("bad input")
				}

			case 4:
				dc.Name = name
				dc.DBNo, err = strconv.Atoi(strDBNo)

			case 5:
				dc.DBByteLength, err = strconv.Atoi(strDBByteLength)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
