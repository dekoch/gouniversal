package appManagement

import (
	"errors"
	"gouniversal/modules/openespm/appConfig"
	"gouniversal/modules/openespm/globalOESPM"
)

func LoadApp(uid string) (appConfig.App, error) {

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.AppConfig.File.Apps); u++ {

		// search app with UUID
		if uid == globalOESPM.AppConfig.File.Apps[u].UUID {

			return globalOESPM.AppConfig.File.Apps[u], nil
		}
	}

	var a appConfig.App
	a.State = -1
	return a, errors.New("LoadApp() app not found")
}

func SaveApp(a appConfig.App) error {

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.AppConfig.File.Apps); u++ {

		// search app with UUID
		if a.UUID == globalOESPM.AppConfig.File.Apps[u].UUID {

			globalOESPM.AppConfig.File.Apps[u] = a
			return nil
		}
	}

	return errors.New("SaveApp() app not found")
}
