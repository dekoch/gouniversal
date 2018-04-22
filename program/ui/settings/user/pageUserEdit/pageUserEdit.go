package pageUserEdit

import (
	"errors"
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/userConfig"
	"gouniversal/program/userManagement"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:User:Edit", page.Lang.Settings.User.UserEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type useredit struct {
		Lang     lang.SettingsUserEdit
		User     userConfig.User
		CmbLang  template.HTML
		CmbState template.HTML
		Groups   template.HTML
	}
	var ue useredit

	ue.Lang = page.Lang.Settings.User.UserEdit

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
	global.UserConfig.Mut.Lock()
	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if id == global.UserConfig.File.User[i].UUID {

			ue.User = global.UserConfig.File.User[i]
		}
	}
	global.UserConfig.Mut.Unlock()

	// combobox Language
	cmbLang := "<select name=\"language\">"

	global.Lang.Mut.Lock()
	for i := 0; i < len(global.Lang.File); i++ {

		cmbLang += "<option value=\"" + global.Lang.File[i].Header.FileName + "\""

		if ue.User.Lang == global.Lang.File[i].Header.FileName {
			cmbLang += " selected"
		}

		cmbLang += ">" + global.Lang.File[i].Header.FileName + "</option>"
	}
	global.Lang.Mut.Unlock()
	cmbLang += "</select>"
	ue.CmbLang = template.HTML(cmbLang)

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

		if ue.User.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	ue.CmbState = template.HTML(cmbState)

	// list of groups
	grouplist := ""

	global.GroupConfig.Mut.Lock()
	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		grouplist += "<tr>"
		grouplist += "<td>" + global.GroupConfig.File.Group[i].Name + "</td>"
		grouplist += "<td><input type=\"checkbox\" name=\"selectedgroups\" value=\"" + global.GroupConfig.File.Group[i].UUID + "\""

		if userManagement.IsUserInGroup(global.GroupConfig.File.Group[i].UUID, ue.User) {

			grouplist += " checked"
		}
		grouplist += "></td></tr>"
	}
	global.GroupConfig.Mut.Unlock()

	ue.Groups = template.HTML(grouplist)

	// display user
	templ, err := template.ParseFiles(global.UiConfig.File.ProgramFileRoot + "settings/useredit.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, ue)
}

func newUser() string {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())

	newuser := make([]userConfig.User, 1)
	newuser[0].UUID = u.String()
	newuser[0].LoginName = u.String()
	newuser[0].Lang = "en"
	newuser[0].State = 1 // active

	global.UserConfig.File.User = append(newuser, global.UserConfig.File.User...)

	global.UserConfig.SaveConfig()

	return u.String()
}

func editUser(r *http.Request, u string) error {

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

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if u == global.UserConfig.File.User[i].UUID {

			iState, err := strconv.Atoi(state)
			if err != nil {
				return err
			}

			selgroups := r.Form["selectedgroups"]

			global.UserConfig.File.User[i].LoginName = loginName
			global.UserConfig.File.User[i].Name = name
			global.UserConfig.File.User[i].State = iState
			global.UserConfig.File.User[i].Lang = sellang
			global.UserConfig.File.User[i].Comment = comment
			global.UserConfig.File.User[i].Groups = selgroups

			pwd := r.FormValue("pwd")

			if pwd != "" {

				if functions.IsEmpty(pwd) {

					global.UserConfig.File.User[i].PWDHash = ""
				} else {

					hash, err := uifunc.HashPassword(pwd)
					if err == nil {
						global.UserConfig.File.User[i].PWDHash = hash
					} else {
						return err
					}
				}
			}

			return global.UserConfig.SaveConfig()
		}
	}

	return errors.New("UUID not found")
}

func deleteUser(u string) error {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	var ul []userConfig.User
	n := make([]userConfig.User, 1)

	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if u != global.UserConfig.File.User[i].UUID {

			n[0] = global.UserConfig.File.User[i]

			ul = append(ul, n...)
		}
	}

	global.UserConfig.File.User = ul

	return global.UserConfig.SaveConfig()
}
