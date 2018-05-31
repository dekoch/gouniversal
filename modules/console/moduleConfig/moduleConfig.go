package moduleConfig

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/console/console"

type ModuleConfigFile struct {
	Header          config.FileHeader
	UIFileRoot      string
	LangFileRoot    string
	RefreshInterval time.Duration
}

type ModuleConfig struct {
	Mut  sync.Mutex
	File ModuleConfigFile
}

func (hc ModuleConfig) SaveConfig() error {

	hc.Mut.Lock()
	defer hc.Mut.Unlock()

	hc.File.Header = config.BuildHeader("console", "console", 1.0, "console config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		hc.File.UIFileRoot = "data/ui/console/1.0/"
		hc.File.LangFileRoot = "data/lang/console/"
		hc.File.RefreshInterval = 500
	}

	b, err := json.Marshal(hc.File)
	if err != nil {
		console.Log(err, "console/moduleConfig.SaveConfig()")
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

	hc.Mut.Lock()
	defer hc.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "console/moduleConfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &hc.File)
	if err != nil {
		console.Log(err, "console/moduleConfig.LoadConfig()")
	}

	if config.CheckHeader(hc.File.Header, "console") == false {
		console.Log("wrong config \""+configFilePath+"\"", "console/moduleConfig.LoadConfig()")
	}

	return err
}
