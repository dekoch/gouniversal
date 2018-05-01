package ui

import (
	"gouniversal/modules/openespm/app"
	"gouniversal/modules/openespm/appManagement"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
	"strings"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesOESPM.Page)
	appPage.Lang = selectLang(nav.User.Lang)

	settings.RegisterPage(appPage, nav)

	// register apps
	globalOESPM.AppConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.AppConfig.File.Apps); i++ {

		a := globalOESPM.AppConfig.File.Apps[i]

		// only active apps
		if a.State == 1 {
			nav.Sitemap.Register("openESPM", "App:openESPM:App:"+a.App+":"+a.UUID, a.Name)
		}
	}
	globalOESPM.AppConfig.Mut.Unlock()
}

func selectLang(l string) langOESPM.File {

	globalOESPM.Lang.Mut.Lock()
	defer globalOESPM.Lang.Mut.Unlock()

	// search lang
	for i := 0; i < len(globalOESPM.Lang.File); i++ {

		if l == globalOESPM.Lang.File[i].Header.FileName {

			return globalOESPM.Lang.File[i]
		}
	}

	// if nothing found
	// search "en"
	for i := 0; i < len(globalOESPM.Lang.File); i++ {

		if "en" == globalOESPM.Lang.File[i].Header.FileName {

			return globalOESPM.Lang.File[i]
		}
	}

	// if nothing found
	// load or create "en"
	return langOESPM.LoadLang("en")
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesOESPM.Page)
	appPage.Content = page.Content
	appPage.Lang = selectLang(nav.User.Lang)

	if nav.IsNext("Settings") {

		settings.Render(appPage, nav, r)

	} else if nav.IsNext("App") {

		// load config for selected app
		err := loadAppConfig(appPage, nav)
		if err == nil {
			app.Render(appPage, nav, r)

			// save config to ram
			appManagement.SaveApp(appPage.App)

			// save config to file
			globalOESPM.AppConfig.Mut.Lock()
			globalOESPM.AppConfig.SaveConfig()
			globalOESPM.AppConfig.Mut.Unlock()
		}
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}

func loadAppConfig(page *typesOESPM.Page, nav *navigation.Navigation) error {

	// search app UUID inside path
	index := strings.LastIndex(nav.Path, ":")

	var u string
	if index > 0 {
		u = nav.Path[index+1:]
	}

	var err error
	page.App, err = appManagement.LoadApp(u)

	return err
}
