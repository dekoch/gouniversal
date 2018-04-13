package appManagement

import (
	"encoding/json"
	"errors"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"

	"github.com/google/uuid"
)

const AppFile = "data/config/openespm/apps"

func SaveConfig(ac typesOESPM.AppConfigFile) error {

	ac.Header = config.BuildHeader("apps", "apps", 1.0, "app config file")

	if _, err := os.Stat(AppFile); os.IsNotExist(err) {
		// if not found, create default file

		newApp := make([]typesOESPM.App, 1)

		u := uuid.Must(uuid.NewRandom())

		newApp[0].UUID = u.String()
		newApp[0].Name = u.String()
		newApp[0].State = 1 // active
		newApp[0].App = "SimpleSwitchV1x0"

		ac.Apps = newApp
	}

	b, err := json.Marshal(ac)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(AppFile, b)

	return err
}

func LoadConfig() typesOESPM.AppConfigFile {

	var ac typesOESPM.AppConfigFile

	if _, err := os.Stat(AppFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveConfig(ac)
	}

	f := new(file.File)
	b, err := f.ReadFile(AppFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(ac.Header, "apps") == false {
		log.Fatal("wrong config")
	}

	return ac
}

func LoadApp(uid string) (typesOESPM.App, error) {

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.AppConfig.File.Apps); u++ {

		// search app with UUID
		if uid == globalOESPM.AppConfig.File.Apps[u].UUID {

			return globalOESPM.AppConfig.File.Apps[u], nil
		}
	}

	var a typesOESPM.App
	a.State = -1
	return a, errors.New("LoadApp() app not found")
}

func SaveApp(a typesOESPM.App) error {

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
