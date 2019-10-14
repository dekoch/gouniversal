package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/hashstor"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/instabackup/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	FileRoot       string
	DBFile         string
	UpdInterv      int // minutes (0=disabled)
	HashReset      int // minutes
	Users          []string
	Hashes         hashstor.HashStor
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "instabackup", ContentName: "instabackup", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/instabackup/1.0/"
	hc.StaticFileRoot = "data/ui/instabackup/1.0/static/"
	hc.LangFileRoot = "data/lang/instabackup/"

	hc.FileRoot = "data/instabackup/"
	hc.DBFile = "data/instabackup/instabackup.db"

	hc.UpdInterv = -1
	hc.HashReset = 5

	hc.Hashes.Add("")
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

func (hc *ModuleConfig) SetUpdInterval(minutes int) {

	mut.Lock()
	defer mut.Unlock()

	hc.UpdInterv = minutes
}

func (hc *ModuleConfig) GetUpdInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.UpdInterv) * time.Minute
}

func (hc *ModuleConfig) SetHashRest(minutes int) {

	mut.Lock()
	defer mut.Unlock()

	hc.HashReset = minutes
}

func (hc *ModuleConfig) GetHashReset() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.HashReset) * time.Minute
}

func (hc *ModuleConfig) AddUser(username string) {

	mut.Lock()
	defer mut.Unlock()

	for i := range hc.Users {

		if hc.Users[i] == username {
			return
		}
	}

	hc.Users = append(hc.Users, username)
}

func (hc *ModuleConfig) GetUserList() []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Users
}
