package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/module/instafollowbot/core/coreconfig"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/types"
	"github.com/google/uuid"
)

const configFilePath = "data/config/instafollowbot/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	FileRoot       string
	DBFile         string
	Cores          []coreconfig.CoreConfig
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "instafollowbot.json", ContentName: "instafollowbot", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/instafollowbot/1.0/"
	hc.StaticFileRoot = "data/ui/instafollowbot/1.0/static/"
	hc.LangFileRoot = "data/lang/instafollowbot/"

	hc.FileRoot = "data/instafollowbot/"
	hc.DBFile = "data/instafollowbot/instafollowbot.db"

	hc.addCore()
}

func (hc ModuleConfig) SaveConfig() error {

	mut.RLock()
	defer mut.RUnlock()

	return hc.saveConfig()
}

func (hc ModuleConfig) saveConfig() error {

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

	mut.Lock()
	defer mut.Unlock()

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		hc.loadDefaults()
		hc.saveConfig()
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

	for i := range hc.Cores {
		hc.Cores[i].LoadConfig()
	}

	return err
}

func (hc *ModuleConfig) Exit(em *types.ExitMessage) error {

	mut.Lock()
	defer mut.Unlock()

	return nil
}

func (hc *ModuleConfig) GetDBFile() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.DBFile
}

func (hc *ModuleConfig) AddCore() (string, error) {

	mut.Lock()
	defer mut.Unlock()

	return hc.addCore()
}

func (hc *ModuleConfig) addCore() (string, error) {

	var co coreconfig.CoreConfig
	co.CoreUUID = uuid.Must(uuid.NewRandom()).String()
	co.LoadDefaults()
	co.LoadConfig()

	hc.Cores = append(hc.Cores, co)

	return co.CoreUUID, nil
}

func (hc *ModuleConfig) GetCoreList() []string {

	mut.RLock()
	defer mut.RUnlock()

	var ret []string

	for i := range hc.Cores {
		ret = append(ret, hc.Cores[i].CoreUUID)
	}

	return ret
}

func (hc *ModuleConfig) SetCoreConfig(co coreconfig.CoreConfig) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range hc.Cores {

		if hc.Cores[i].CoreUUID == co.CoreUUID {

			hc.Cores[i] = co
			return nil
		}
	}

	return errors.New("uuid not found")
}

func (hc *ModuleConfig) GetCoreConfig(uid string) (coreconfig.CoreConfig, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := range hc.Cores {

		if hc.Cores[i].CoreUUID == uid {

			return hc.Cores[i], nil
		}
	}

	var co coreconfig.CoreConfig
	return co, errors.New("uuid not found")
}
