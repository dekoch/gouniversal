package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/station"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/gasprice/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	FileRoot       string
	DBFile         string
	UpdInterv      int // minutes (0=disabled)
	SaveToDB       bool
	SaveToCSV      bool
	LoadFromDB     bool
	LoadFromCSV    bool
	ImportCSVtoDB  bool
	GasTypes       []string
	Stations       station.StationList
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "gasprice", ContentName: "gasprice", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/gasprice/1.0/"
	hc.StaticFileRoot = "data/ui/gasprice/1.0/static/"
	hc.LangFileRoot = "data/lang/gasprice/"
	hc.FileRoot = "data/gasprice/"
	hc.DBFile = "data/gasprice/gasprice.db"
	hc.UpdInterv = 0
	hc.SaveToDB = false
	hc.SaveToCSV = true
	hc.LoadFromDB = false
	hc.LoadFromCSV = true
	hc.ImportCSVtoDB = false

	gasTypes := make([]string, 3)
	gasTypes[0] = "Diesel"
	gasTypes[1] = "Super"
	gasTypes[2] = "Super E10"
	hc.GasTypes = gasTypes
}

func (hc ModuleConfig) SaveConfig() error {

	hc.Header = config.BuildHeaderWithStruct(header)

	b, err := json.Marshal(hc)
	if err != nil {
		console.Log(err, "")
		return err
	}

	err = file.WriteFile(configFilePath+header.FileName, b)
	if err != nil {
		console.Log(err, "")
	}

	return err
}

func (hc *ModuleConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		hc.loadDefaults()
		hc.SaveConfig()
	}

	b, err := file.ReadFile(configFilePath + header.FileName)
	if err != nil {
		console.Log(err, "")
		hc.loadDefaults()
	} else {
		err = json.Unmarshal(b, &hc)
		if err != nil {
			console.Log(err, "")
			hc.loadDefaults()
		}
	}

	if config.CheckHeader(hc.Header, header.ContentName) == false {
		err = errors.New("wrong config \"" + configFilePath + header.FileName + "\"")
		console.Log(err, "")
		hc.loadDefaults()
	}

	return err
}

func (hc *ModuleConfig) GetUpdInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.UpdInterv) * time.Minute
}

func (hc *ModuleConfig) GetGasTypes() []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.GasTypes
}
