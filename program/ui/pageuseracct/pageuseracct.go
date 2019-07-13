package pageuseracct

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	if nav.User.State == userconfig.StateActive {
		nav.Sitemap.Register("Account", "Account:UserAccount", page.Lang.UserAccount.Title)
	}
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	var err error

	type content struct {
		Lang    lang.UserAccount
		User    userconfig.User
		CmbLang template.HTML
	}
	var c content

	c.Lang = page.Lang.UserAccount

	// Form input
	switch r.FormValue("edit") {
	case "apply":
		var user userconfig.User

		user, err = editUser(r, nav.User.UUID, page, nav)
		if err == nil {
			nav.User = user
			nav.RedirectPath("Account:User", false)
			return
		}
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	// copy user from array
	c.User, err = global.UserConfig.Get(nav.User.UUID)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

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

	// display user
	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"useracct.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func editUser(r *http.Request, uid string, page *types.Page, nav *navigation.Navigation) (userconfig.User, error) {

	var (
		err             error
		loginName       string
		name            string
		selLang         string
		setNewPassword  bool
		newPassword     string
		confirmPassword string
		u               userconfig.User
	)

	func() {

		for i := 0; i <= 9; i++ {

			switch i {
			case 0:
				loginName, err = functions.CheckFormInput("loginname", r)

			case 1:
				name, err = functions.CheckFormInput("name", r)

			case 2:
				selLang = r.FormValue("language")

			case 3:
				// check input
				if functions.IsEmpty(loginName) ||
					functions.IsEmpty(selLang) {

					err = errors.New("bad input")
				}

			case 4:
				u, err = global.UserConfig.Get(uid)

			case 5:
				u.LoginName = loginName
				u.Name = name
				u.Lang = selLang

			case 6:
				newPassword = r.FormValue("newpwd")
				confirmPassword = r.FormValue("confirmpwd")

				if functions.IsEmpty(newPassword) == false {

					setNewPassword = true

					if newPassword != confirmPassword {
						err = errors.New(page.Lang.UserAccount.PasswordMismatch)
					}
				}

			case 7:
				if setNewPassword {

					oldPassword := r.FormValue("oldpwd")

					if uifunc.CheckLogin(nav.User.LoginName, oldPassword) {

						hash, err := uifunc.HashPassword(newPassword)
						if err == nil {
							u.PWDHash = hash
						}
					} else {
						err = errors.New(page.Lang.UserAccount.IncorrectPassword)
					}
				}

			case 8:
				err = global.UserConfig.Edit(u)

			case 9:
				err = global.UserConfig.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return u, err
}
