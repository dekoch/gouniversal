package pageuseredit

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/program/usermanagement"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"

	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:User:Edit", page.Lang.Settings.User.UserEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	var err error

	type content struct {
		Lang     lang.SettingsUserEdit
		User     userconfig.User
		CmbLang  template.HTML
		CmbState template.HTML
		Groups   template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.User.UserEdit

	// Form input
	id := nav.Parameter("UUID")

	if id == "new" {

		id, err = newUser()
		if err == nil {
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
			return
		}
	}

	switch r.FormValue("edit") {
	case "apply":
		err = editUser(r, id)
		if err == nil {
			nav.RedirectPath("Program:Settings:User:List", false)
			return
		}

	case "delete":
		err = deleteUser(id)
		if err == nil {
			nav.RedirectPath("Program:Settings:User:List", false)
			return
		}
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	// copy user from array
	c.User, err = global.UserConfig.Get(id)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	// combobox Language
	cmbLang := "<select name=\"language\">"

	err = global.Lang.LoadLangFiles()
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	langFiles := global.Lang.ListNames()

	for i := 0; i < len(langFiles); i++ {

		cmbLang += "<option value=\"" + langFiles[i] + "\""

		if c.User.Lang == langFiles[i] {
			cmbLang += " selected"
		}

		cmbLang += ">" + langFiles[i] + "</option>"
	}

	cmbLang += "</select>"
	c.CmbLang = template.HTML(cmbLang)

	// combobox State
	cmbState := "<select name=\"state\">"
	statetext := ""

	var us userconfig.UserState

	for i := 0; i <= 2; i++ {

		us = userconfig.UserState(i)

		switch us {
		case userconfig.StatePublic:
			statetext = page.Lang.Settings.User.UserEdit.States.Public
		case userconfig.StateActive:
			statetext = page.Lang.Settings.User.UserEdit.States.Active
		case userconfig.StateInactive:
			statetext = page.Lang.Settings.User.UserEdit.States.Inactive
		}

		cmbState += "<option value=\"" + strconv.Itoa(i) + "\""

		if c.User.State == us {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	c.CmbState = template.HTML(cmbState)

	// list of groups
	grouplist := ""

	groups := global.GroupConfig.List()

	for _, g := range groups {

		grouplist += "<tr>"
		grouplist += "<td><a href=\"/?path=Program:Settings:Group:Edit$UUID=" + g.UUID + "\">" + g.Name + "</a></td>"
		grouplist += "<td><input type=\"checkbox\" name=\"selectedgroups\" value=\"" + g.UUID + "\""

		if usermanagement.IsUserInGroup(g.UUID, c.User) {

			grouplist += " checked"
		}
		grouplist += "></td></tr>"
	}

	c.Groups = template.HTML(grouplist)

	// display user
	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"settings/useredit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newUser() (string, error) {

	u := uuid.Must(uuid.NewRandom())

	var newUser userconfig.User
	newUser.UUID = u.String()
	newUser.LoginName = u.String()
	newUser.Lang = "en"
	newUser.State = 1 // active

	global.UserConfig.Add(newUser)

	err := global.UserConfig.SaveConfig()

	return u.String(), err
}

func editUser(r *http.Request, uid string) error {

	var (
		err       error
		loginName string
		name      string
		strState  string
		intState  int
		selLang   string
		comment   string
		u         userconfig.User
	)

	func() {

		for i := 0; i <= 12; i++ {

			switch i {
			case 0:
				loginName, err = functions.CheckFormInput("loginname", r)

			case 1:
				name, err = functions.CheckFormInput("name", r)

			case 2:
				strState, err = functions.CheckFormInput("state", r)

			case 3:
				selLang = r.FormValue("language")

			case 4:
				comment, err = functions.CheckFormInput("comment", r)

			case 5:
				// check input
				if functions.IsEmpty(loginName) ||
					functions.IsEmpty(strState) ||
					functions.IsEmpty(selLang) {

					err = errors.New("bad input")
				}

			case 6:
				intState, err = strconv.Atoi(strState)

			case 7:
				if intState < 0 ||
					intState > 2 {

					err = errors.New("bad input")
				}

			case 8:
				u, err = global.UserConfig.Get(uid)

			case 9:
				u.LoginName = loginName
				u.Name = name
				u.State = userconfig.UserState(intState)
				u.Lang = selLang
				u.Comment = comment
				u.Groups = r.Form["selectedgroups"]

			case 10:
				pwd := r.FormValue("pwd")

				if functions.IsEmpty(pwd) == false {

					hash, err := uifunc.HashPassword(pwd)
					if err == nil {
						u.PWDHash = hash
					}
				}

			case 11:
				err = global.UserConfig.Edit(u)

			case 12:
				err = global.UserConfig.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteUser(uid string) error {

	global.UserConfig.Delete(uid)

	return global.UserConfig.SaveConfig()
}
