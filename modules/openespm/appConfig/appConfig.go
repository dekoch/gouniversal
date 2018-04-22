package appConfig

import (
	"encoding/json"
	"fmt"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

const configFilePath = "data/config/openespm/apps"

type App struct {
	UUID    string
	Name    string
	State   int
	Comment string
	App     string
	Config  string
}

func (a App) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(a.Config), &v)
}

func (a *App) Marshal(v interface{}) error {
	b, err := json.Marshal(v)
	if err == nil {
		a.Config = string(b[:])
	} else {
		fmt.Println(err)
	}

	return err
}

type AppConfigFile struct {
	Header config.FileHeader
	Apps   []App
}

type AppConfig struct {
	Mut  sync.Mutex
	File AppConfigFile
}

func (ac AppConfig) SaveConfig() error {

	ac.File.Header = config.BuildHeader("apps", "apps", 1.0, "app config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newApp := make([]App, 1)

		u := uuid.Must(uuid.NewRandom())

		newApp[0].UUID = u.String()
		newApp[0].Name = u.String()
		newApp[0].State = 1 // active
		newApp[0].App = "SimpleSwitchV1x0"

		ac.File.Apps = newApp
	}

	b, err := json.Marshal(ac.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (ac *AppConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		ac.SaveConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &ac.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(ac.File.Header, "apps") == false {
		log.Fatal("wrong config")
	}

	return err
}
