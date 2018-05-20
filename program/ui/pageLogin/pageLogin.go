package pageLogin

import (
	"net/http"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Account", "Account:Login", page.Lang.Menu.Account.Login)
	nav.Sitemap.Register("Account", "Account:Logout", page.Lang.Menu.Account.Logout)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang lang.Login
	}
	var c content

	c.Lang = page.Lang.Login

	p, err := functions.PageToString(global.UiConfig.File.ProgramFileRoot+"login.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
