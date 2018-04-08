package ui

import (
	"gouniversal/modules/openespm/app"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/modules/openespm/ui/settings"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	appPage := new(typesOESPM.Page)
	appPage.Lang = selectLang(nav.User.Lang)

	settings.RegisterPage(appPage, nav)

	// register devices
	globalOESPM.DeviceConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.DeviceConfig.File.Devices); i++ {

		dev := globalOESPM.DeviceConfig.File.Devices[i]

		// only active devices
		if dev.State == 1 {
			nav.Sitemap.Register("App", "App:openESPM:App:"+dev.App+":"+dev.UUID, dev.Name)
		}
	}
	globalOESPM.DeviceConfig.Mut.Unlock()
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

		app.Render(appPage, nav, r)
	}

	page.Content += appPage.Content
}
