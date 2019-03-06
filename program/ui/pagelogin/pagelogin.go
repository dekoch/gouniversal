package pagelogin

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
	"github.com/dekoch/gouniversal/shared/clientinfo"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/timeout"
	token "github.com/dekoch/gouniversal/shared/token"
	"github.com/dekoch/gouniversal/shared/types"
)

var (
	tokens token.Token
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Account", "Account:Login", page.Lang.Menu.Account.Login)
	nav.Sitemap.Register("Account", "Account:Logout", page.Lang.Menu.Account.Logout)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang  lang.Login
		UUID  template.HTML
		Token template.HTML
	}
	var c content

	c.Lang = page.Lang.Login

	if r.FormValue("edit") == "login" {

		var (
			err   error
			id    string
			tok   string
			name  string
			delay timeout.TimeOut
		)

		valid := true

		delay.Start(3000)

		for i := 0; i <= 6; i++ {

			if valid == false {
				continue
			}

			switch i {
			case 0:
				id = r.FormValue("uuid")

				if functions.IsEmpty(id) {
					valid = false
				}

			case 1:
				tok = r.FormValue("token")

				if functions.IsEmpty(tok) {
					valid = false
				}

			case 2:
				if tokens.Check(id, tok) == false {
					valid = false
				}

			case 3:
				name, err = functions.CheckFormInput("name", r)
				if err != nil {
					valid = false
				}

			case 4:
				pwd := r.FormValue("pwd")

				if uifunc.CheckLogin(name, pwd) == false {
					valid = false
				}

			case 5:
				if global.Guests.MaxLoginAttempts(nav.User.UUID, global.UIConfig.MaxGuests) {
					valid = false
				}

			case 6:
				// load user
				nav.User, err = global.UserConfig.Get(uifunc.LoginNameToUUID(name))
				if err != nil {
					valid = false
				}
			}
		}

		if valid {

			console.Log("\""+name+"\" logged in ("+clientinfo.String(r)+")", "Login")

			nav.RedirectPath("Program:Home", false)
		} else {

			if functions.IsEmpty(name) == false {
				console.Log("failed login with user \""+name+"\" ("+clientinfo.String(r)+")", "Login")
			}

			<-delay.ElapsedChan()
		}
	}

	c.UUID = template.HTML(nav.User.UUID)

	tokens.SetMaxTokens(global.UserConfig.GetUserCnt() + global.UIConfig.MaxGuests)
	c.Token = template.HTML(tokens.New(nav.User.UUID))

	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"login.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
