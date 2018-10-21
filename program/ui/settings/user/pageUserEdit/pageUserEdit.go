package pageUserEdit

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
	"github.com/dekoch/gouniversal/program/userConfig"
	"github.com/dekoch/gouniversal/program/userManagement"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:User:Edit", page.Lang.Settings.User.UserEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang     lang.SettingsUserEdit
		User     userConfig.User
		CmbLang  template.HTML
		CmbState template.HTML
		Groups   template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.User.UserEdit

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newUser()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		err := editUser(r, id)
		if err == nil {
			nav.RedirectPath("Program:Settings:User:List", false)
		}

	} else if button == "delete" {

		err := deleteUser(id)
		if err == nil {
			nav.RedirectPath("Program:Settings:User:List", false)
		}
	}

	// copy user from array
	var err error
	c.User, err = global.UserConfig.Get(id)

	// combobox Language
	cmbLang := "<select name=\"language\">"

	global.Lang.LoadLangFiles()

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

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			statetext = page.Lang.Settings.User.UserEdit.States.Public
		case 1:
			statetext = page.Lang.Settings.User.UserEdit.States.Active
		case 2:
			statetext = page.Lang.Settings.User.UserEdit.States.Inactive
		}

		cmbState += "<option value=\"" + strconv.Itoa(i) + "\""

		if c.User.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	c.CmbState = template.HTML(cmbState)

	// list of groups
	grouplist := ""

	groups := global.GroupConfig.List()

	for i := 0; i < len(groups); i++ {

		g := groups[i]

		grouplist += "<tr>"
		grouplist += "<td>" + g.Name + "</td>"
		grouplist += "<td><input type=\"checkbox\" name=\"selectedgroups\" value=\"" + g.UUID + "\""

		if userManagement.IsUserInGroup(g.UUID, c.User) {

			grouplist += " checked"
		}
		grouplist += "></td></tr>"
	}

	c.Groups = template.HTML(grouplist)

	// display user
	p, err := functions.PageToString(global.UiConfig.ProgramFileRoot+"settings/useredit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newUser() string {

	u := uuid.Must(uuid.NewRandom())

	var newUser userConfig.User
	newUser.UUID = u.String()
	newUser.LoginName = u.String()
	newUser.Lang = "en"
	newUser.State = 1 // active

	global.UserConfig.Add(newUser)

	global.UserConfig.SaveConfig()

	return u.String()
}

func editUser(r *http.Request, uid string) error {

	loginName, _ := functions.CheckFormInput("loginname", r)
	name, errName := functions.CheckFormInput("name", r)
	state, _ := functions.CheckFormInput("state", r)
	sellang := r.FormValue("language")
	comment, errComment := functions.CheckFormInput("comment", r)

	// check input
	if functions.IsEmpty(loginName) ||
		functions.IsEmpty(state) ||
		govalidator.IsNumeric(state) == false ||
		functions.IsEmpty(sellang) ||
		// content not required
		errName != nil ||
		errComment != nil {

		return errors.New("bad input")
	}

	iState, err := strconv.Atoi(state)
	if err != nil {
		return err
	}

	selgroups := r.Form["selectedgroups"]

	u, err := global.UserConfig.Get(uid)
	if err != nil {
		return err
	}

	u.LoginName = loginName
	u.Name = name
	u.State = iState
	u.Lang = sellang
	u.Comment = comment
	u.Groups = selgroups

	pwd := r.FormValue("pwd")

	if pwd != "" {

		if functions.IsEmpty(pwd) {

			u.PWDHash = ""
		} else {

			hash, err := uifunc.HashPassword(pwd)
			if err == nil {
				u.PWDHash = hash
			} else {
				return err
			}
		}
	}

	err = global.UserConfig.Edit(u)
	if err != nil {
		return err
	}

	return global.UserConfig.SaveConfig()
}

func deleteUser(uid string) error {

	global.UserConfig.Delete(uid)

	return global.UserConfig.SaveConfig()
}
