package moduleConfig

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
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
		console.Log(err, "fileshare/moduleConfig.SaveConfig()")
	}

	err = file.WriteFile(configFilePath, b)

	return err
}

func (hc *ModuleConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		hc.SaveConfig()
	}

	b, err := file.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "fileshare/moduleConfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &hc.File)
	if err != nil {
		console.Log(err, "fileshare/moduleConfig.LoadConfig()")
	}

	if config.CheckHeader(hc.File.Header, "fileshare") == false {
		console.Log("wrong config \""+configFilePath+"\"", "fileshare/moduleConfig.LoadConfig()")
	}

	return err
}
