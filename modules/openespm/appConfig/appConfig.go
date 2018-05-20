package appConfig

import (
	"encoding/json"
	"errors"
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

func (c *AppConfig) SaveConfig() error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	c.File.Header = config.BuildHeader("apps", "apps", 1.0, "app config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newApp := make([]App, 1)

		u := uuid.Must(uuid.NewRandom())

		newApp[0].UUID = u.String()
		newApp[0].Name = u.String()
		newApp[0].State = 1 // active
		newApp[0].App = "SimpleSwitchV1x0"

		c.File.Apps = newApp
	}

	b, err := json.Marshal(c.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (c *AppConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.SaveConfig()
	}

	c.Mut.Lock()
	defer c.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &c.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(c.File.Header, "apps") == false {
		log.Fatal("wrong config")
	}

	return err
}

func (c *AppConfig) Add(a App) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	newApp := make([]App, 1)

	newApp[0] = a

	c.File.Apps = append(c.File.Apps, newApp...)
}

func (c *AppConfig) Edit(uid string, a App) error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.Apps); i++ {

		if uid == c.File.Apps[i].UUID {

			c.File.Apps[i] = a
			return nil
		}
	}

	return errors.New("Edit() app \"" + a.UUID + "\" not found")
}

func (c *AppConfig) Get(uid string) (App, error) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.Apps); i++ {

		if uid == c.File.Apps[i].UUID {

			return c.File.Apps[i], nil
		}
	}

	var a App
	a.State = -1
	return a, errors.New("Get() app \"" + uid + "\" not found")
}

func (c *AppConfig) List() []App {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	return c.File.Apps
}

func (c *AppConfig) Delete(uid string) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	var al []App
	n := make([]App, 1)

	for i := 0; i < len(c.File.Apps); i++ {

		if uid != c.File.Apps[i].UUID {

			n[0] = c.File.Apps[i]

			al = append(al, n...)
		}
	}

	c.File.Apps = al
}
