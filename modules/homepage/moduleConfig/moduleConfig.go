package moduleConfig

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"
)

const configFilePath = "data/config/homepage/homepage"

type ModuleConfigFile struct {
	Header     config.FileHeader
	UIFileRoot string
}

type ModuleConfig struct {
	Mut  sync.Mutex
	File ModuleConfigFile
}

func (hc ModuleConfig) SaveConfig() error {

	hc.File.Header = config.BuildHeader("homepage", "homepage", 1.0, "homepage config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		hc.File.UIFileRoot = "data/ui/homepage/1.0/"
	}

	b, err := json.Marshal(hc.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (hc *ModuleConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		hc.SaveConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &hc.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(hc.File.Header, "homepage") == false {
		log.Fatal("wrong config")
	}

	return err
}
