package moduleConfig

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
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

		hc.File.UIFileRoot = "data/homepage/"
	}

	b, err := json.Marshal(hc.File)
	if err != nil {
		console.Log(err, "homepage/moduleConfig.SaveConfig()")
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
		console.Log(err, "homepage/moduleConfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &hc.File)
	if err != nil {
		console.Log(err, "homepage/moduleConfig.LoadConfig()")
	}

	if config.CheckHeader(hc.File.Header, "homepage") == false {
		console.Log("wrong config \""+configFilePath+"\"", "homepage/moduleConfig.LoadConfig()")
	}

	return err
}
