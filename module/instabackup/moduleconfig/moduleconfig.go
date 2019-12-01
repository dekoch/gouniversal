package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/instabackup/userconfig"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/hashstor"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/types"
)

const configFilePath = "data/config/instabackup/"

type ModuleConfig struct {
	Header          config.FileHeader
	UIFileRoot      string
	StaticFileRoot  string
	LangFileRoot    string
	FileRoot        string
	DBFile          string
	CheckInterv     int // minutes (0=disabled)
	UpdInterv       int // minutes
	HashReset       int // minutes
	MaxClientTokens int
	Users           []userconfig.UserConfig
	Hashes          hashstor.HashStor
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

	hc.CheckInterv = 5
	hc.UpdInterv = 30
	hc.HashReset = 5
	hc.MaxClientTokens = 5

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

	hc.Hashes.Init()

	return err
}

func (hc *ModuleConfig) Exit(em *types.ExitMessage) error {

	mut.Lock()
	defer mut.Unlock()

	var l []userconfig.UserConfig

	for iu := range hc.Users {

		found := false

		for i := range em.Users {

			if hc.Users[iu].User == em.Users[i] {
				found = true
			}
		}

		if found {
			l = append(l, hc.Users[iu])
		}
	}

	hc.Users = l

	return nil
}

func (hc *ModuleConfig) selectUser(user string) *userconfig.UserConfig {

	for iu := range hc.Users {

		if hc.Users[iu].User == user {
			return &hc.Users[iu]
		}
	}
	// create new user
	var n userconfig.UserConfig
	n.LoadDefaults(user)

	hc.Users = append(hc.Users, n)
	// return new user from array
	for iu := range hc.Users {

		if hc.Users[iu].User == user {
			return &hc.Users[iu]
		}
	}

	return nil
}

func (hc *ModuleConfig) SetCheckInterval(minutes int) {

	mut.Lock()
	defer mut.Unlock()

	hc.CheckInterv = minutes
}

func (hc *ModuleConfig) GetCheckInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.CheckInterv) * time.Minute
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

func (hc *ModuleConfig) SetMaxClientTokens(n int) {

	mut.Lock()
	defer mut.Unlock()

	hc.MaxClientTokens = n
}

func (hc *ModuleConfig) GetMaxClientTokens() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.MaxClientTokens
}

func (hc *ModuleConfig) GetAllIDs() []string {

	mut.RLock()
	defer mut.RUnlock()

	var ret []string

	for i := range hc.Users {
		ret = append(ret, hc.Users[i].GetIDList()...)
	}

	return ret
}

func (hc *ModuleConfig) AddIDToUser(user, instaid string) {

	mut.Lock()
	defer mut.Unlock()

	n := hc.selectUser(user)
	n.AddID(instaid)
}

func (hc *ModuleConfig) GetIDFromUser(user string) []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.selectUser(user).GetIDList()
}
