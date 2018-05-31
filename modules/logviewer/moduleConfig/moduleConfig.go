package moduleConfig

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/logviewer/logviewer"

type ModuleConfigFile struct {
	Header       config.FileHeader
	UIFileRoot   string
	LangFileRoot string
	LogFileRoot  string
}

type ModuleConfig struct {
	Mut  sync.Mutex
	File ModuleConfigFile
}

func (hc ModuleConfig) SaveConfig() error {

	hc.Mut.Lock()
	defer hc.Mut.Unlock()

	hc.File.Header = config.BuildHeader("logviewer", "logviewer", 1.0, "logviewer config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		hc.File.UIFileRoot = "data/ui/logviewer/1.0/"
		hc.File.LangFileRoot = "data/lang/logviewer/"
		hc.File.LogFileRoot = "data/log/"
	}

	b, err := json.Marshal(hc.File)
	if err != nil {
		console.Log(err, "logviewer/moduleConfig.SaveConfig()")
	}

	err = file.WriteFile(configFilePath, b)

	return err
}

func (hc *ModuleConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		hc.SaveConfig()
	}

	hc.Mut.Lock()
	defer hc.Mut.Unlock()

	b, err := file.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "logviewer/moduleConfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &hc.File)
	if err != nil {
		console.Log(err, "logviewer/moduleConfig.LoadConfig()")
	}

	if config.CheckHeader(hc.File.Header, "logviewer") == false {
		console.Log("wrong config \""+configFilePath+"\"", "logviewer/moduleConfig.LoadConfig()")
	}

	return err
}
