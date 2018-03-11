package pageUserEdit

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/types"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"gouniversal/program/userManagement"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:User:Edit", page.Lang.Settings.User.UserEdit.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	type useredit struct {
		Lang    lang.SettingsUserEdit
		User    types.User
		CmbLang template.HTML
		Groups  template.HTML
	}
	var ue useredit

	ue.Lang = page.Lang.Settings.User.UserEdit

	button := r.FormValue("edit")

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newUser()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		editUser(r, id)

		nav.RedirectPath("Program:Settings:User:List", false)

	} else if button == "delete" {

		deleteUser(id)

		nav.RedirectPath("Program:Settings:User:List", false)
	}

	// copy user from array
	global.UserConfig.Mut.Lock()
	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if id == global.UserConfig.File.User[i].UUID {

			ue.User = global.UserConfig.File.User[i]
		}
	}
	global.UserConfig.Mut.Unlock()

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

	cmbLang := "<select name=\"language\">"

	uiglobal.Lang.Mut.Lock()
	for i := 0; i < len(uiglobal.Lang.File); i++ {

		cmbLang += "<option value=\"" + uiglobal.Lang.File[i].Header.FileName + "\""

		if ue.User.Lang == uiglobal.Lang.File[i].Header.FileName {
			cmbLang += " selected"
		}

		cmbLang += ">" + uiglobal.Lang.File[i].Header.FileName + "</option>"
	}
	uiglobal.Lang.Mut.Unlock()
	cmbLang += "</select>"
	ue.CmbLang = template.HTML(cmbLang)

	// display user
	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/settings/useredit.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += uifunc.TemplToString(templ, ue)
}

func newUser() string {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())

	newuser := make([]types.User, 1)
	newuser[0].UUID = u.String()
	newuser[0].LoginName = u.String()
	newuser[0].Lang = "en"

	global.UserConfig.File.User = append(newuser, global.UserConfig.File.User...)

	userManagement.SaveUser(global.UserConfig.File)

	return u.String()
}

func editUser(r *http.Request, u string) {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	loginName := uifunc.CheckFormInput("loginname", r)
	state := uifunc.CheckFormInput("state", r)

	if uifunc.CheckInput(loginName, uifunc.STRING) &&
		uifunc.CheckInput(state, uifunc.INT) {

		for i := 0; i < len(global.UserConfig.File.User); i++ {

			if u == global.UserConfig.File.User[i].UUID {

				iState, err := strconv.Atoi(state)

				if err == nil {
					name := uifunc.CheckFormInput("name", r)
					sellang := r.FormValue("language")
					selgroups := r.Form["selectedgroups"]

					global.UserConfig.File.User[i].LoginName = loginName
					global.UserConfig.File.User[i].Name = name
					global.UserConfig.File.User[i].State = iState
					global.UserConfig.File.User[i].Lang = sellang
					global.UserConfig.File.User[i].Groups = selgroups

					pwd := r.FormValue("pwd")

					if uifunc.CheckInput(pwd, uifunc.STRING) {

						hash, err := uifunc.HashPassword(pwd)

						if err == nil {
							global.UserConfig.File.User[i].PWDHash = hash
						}
					}
				}

				userManagement.SaveUser(global.UserConfig.File)
			}
		}
	}
}

func deleteUser(u string) {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	var ul []types.User
	n := make([]types.User, 1)

	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if u != global.UserConfig.File.User[i].UUID {

			n[0] = global.UserConfig.File.User[i]

			ul = append(ul, n...)
		}
	}

	global.UserConfig.File.User = ul

	userManagement.SaveUser(global.UserConfig.File)
}
