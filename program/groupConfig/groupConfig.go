package groupConfig

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"
)

const configFilePath = "data/config/group"

type Group struct {
	UUID         string
	Name         string
	State        int
	Comment      string
	CanModify    bool
	AllowedPages []string
}

type GroupConfigFile struct {
	Header config.FileHeader
	Group  []Group
}

type GroupConfig struct {
	Mut  sync.Mutex
	File GroupConfigFile
}

func (gc GroupConfig) SaveConfig() error {

	gc.File.Header = config.BuildHeader("group", "groups", 1.0, "group config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newgroup := make([]Group, 1)

		newgroup[0].UUID = "admin"
		newgroup[0].Name = "admin"
		newgroup[0].State = 1 // active

		pages := []string{"Program:Settings:User", "Program:Settings:User:List", "Program:Settings:User:Edit", "Program:Settings:Group", "Program:Settings:Group:List", "Program:Settings:Group:Edit"}
		newgroup[0].AllowedPages = pages

		gc.File.Group = newgroup
	}

	b, err := json.Marshal(gc.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (gc *GroupConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		gc.SaveConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &gc.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(gc.File.Header, "groups") == false {
		log.Fatal("wrong config \"" + configFilePath + "\"")
	}

	return err
}
