package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/module/backup/backupitem"
	"github.com/dekoch/gouniversal/module/backup/userconfig"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/backup/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	FileRoot       string
	Users          []userconfig.UserConfig
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "backup", ContentName: "backup", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/backup/1.0/"
	hc.StaticFileRoot = "data/ui/backup/1.0/static/"
	hc.LangFileRoot = "data/lang/backup/"

	hc.FileRoot = "data/backup/"

	hc.Users = []userconfig.UserConfig{}
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

	return err
}

func (hc *ModuleConfig) selectUser(user string) *userconfig.UserConfig {

	for iu := range hc.Users {

		if hc.Users[iu].User == user {
			return &hc.Users[iu]
		}
	}
	// create new user
	var n userconfig.UserConfig
	n.LoadDefaults(user, hc.FileRoot)

	hc.Users = append(hc.Users, n)
	// return new user from array
	for iu := range hc.Users {

		if hc.Users[iu].User == user {
			return &hc.Users[iu]
		}
	}

	return nil
}

func (hc *ModuleConfig) AddItem(user string, item backupitem.BackupItem) {

	mut.Lock()
	defer mut.Unlock()

	n := hc.selectUser(user)
	n.AddItem(item)
}

func (hc *ModuleConfig) GetItems(user string) []backupitem.BackupItem {

	mut.RLock()
	defer mut.RUnlock()

	return hc.selectUser(user).GetItems()
}

func (hc *ModuleConfig) AddExclude(user string, path string) {

	mut.Lock()
	defer mut.Unlock()

	n := hc.selectUser(user)
	n.AddExclude(path)
}

func (hc *ModuleConfig) GetExclude(user string) []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.selectUser(user).GetExclude()
}
