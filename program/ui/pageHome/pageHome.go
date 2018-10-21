package pageHome

import (
	"net/http"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "Program:Home", page.Lang.Home.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Temp string
	}
	var c content

	p, err := functions.PageToString(global.UiConfig.ProgramFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
