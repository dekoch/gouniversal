package moduleConfig

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/fileshare/fileshare"

type ModuleConfigFile struct {
	Header     config.FileHeader
	UIFileRoot string
	FileRoot   string
}

type ModuleConfig struct {
	Mut  sync.Mutex
	File ModuleConfigFile
}

func (hc ModuleConfig) SaveConfig() error {

	hc.File.Header = config.BuildHeader("fileshare", "fileshare", 1.0, "fileshare config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		hc.File.UIFileRoot = "data/ui/fileshare/1.0/"
		hc.File.FileRoot = "data/fileshare/"
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

	if config.CheckHeader(hc.File.Header, "fileshare") == false {
		log.Fatal("wrong config")
	}

	return err
}
