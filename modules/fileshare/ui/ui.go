package ui

import (
	"gouniversal/modules/fileshare/global"
	"gouniversal/modules/fileshare/lang"
	"gouniversal/modules/fileshare/typesFileshare"
	"gouniversal/modules/fileshare/ui/home"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesFileshare.Page)
	appPage.Lang = selectLang(nav.User.Lang)

	home.RegisterPage(appPage, nav)
}

func selectLang(l string) lang.File {

	global.Lang.Mut.Lock()
	defer global.Lang.Mut.Unlock()

	// search lang
	for i := 0; i < len(global.Lang.Files); i++ {

		if l == global.Lang.Files[i].Header.FileName {

			return global.Lang.Files[i]
		}
	}

	// if nothing found
	// search "en"
	for i := 0; i < len(global.Lang.Files); i++ {

		if "en" == global.Lang.Files[i].Header.FileName {

			return global.Lang.Files[i]
		}
	}

	// if nothing found
	// load or create "en"
	return lang.LoadLang("en")
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesFileshare.Page)
	appPage.Content = page.Content
	appPage.Lang = selectLang(nav.User.Lang)

	if nav.IsNext("Home") {

		home.Render(appPage, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}

func LoadConfig() {

	home.LoadConfig()
}
