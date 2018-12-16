package ui

import (
	"net/http"
	"strings"

	"github.com/dekoch/gouniversal/module/openespm/app"
	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/module/openespm/ui/settings"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesOESPM.Page)
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	settings.RegisterPage(appPage, nav)

	// register apps
	apps := global.AppConfig.List()
	for i := 0; i < len(apps); i++ {

		a := apps[i]

		// only active apps
		if a.State == 1 {
			nav.Sitemap.Register("openESPM", "App:openESPM:App:"+a.App+":"+a.UUID, a.Name)
		}
	}
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	appPage := new(typesOESPM.Page)
	appPage.Content = page.Content
	global.Lang.SelectLang(nav.User.Lang, &appPage.Lang)

	if nav.IsNext("Settings") {

		settings.Render(appPage, nav, r)

	} else if nav.IsNext("App") {

		// load config for selected app
		err := loadAppConfig(appPage, nav)
		if err == nil {
			app.Render(appPage, nav, r)

			// save config to ram
			global.AppConfig.Edit(appPage.App.UUID, appPage.App)

			// save config to file
			global.AppConfig.SaveConfig()
		}
	} else {
		nav.RedirectPath("404", true)
	}

	page.Content += appPage.Content
}

func loadAppConfig(page *typesOESPM.Page, nav *navigation.Navigation) error {

	// search app UUID inside path
	index := strings.LastIndex(nav.Path, ":")

	var uid string
	if index > 0 {
		uid = nav.Path[index+1:]
	}

	var err error
	page.App, err = global.AppConfig.Get(uid)

	return err
}
