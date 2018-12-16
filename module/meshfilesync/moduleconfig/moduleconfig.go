package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/module/meshfilesync/syncfile"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/meshfilesync/"

type ModuleConfig struct {
	Header       config.FileHeader
	UIFileRoot   string
	LangFileRoot string
	FileRoot     string
	TempRoot     string
	MaxFileSize  float64 // Megabytes
	AutoAdd      bool
	AutoUpdate   bool
	AutoDelete   bool
	LocalFiles   []syncfile.SyncFile
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "meshfilesync", ContentName: "meshfilesync", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/meshfilesync/1.0/"
	hc.LangFileRoot = "data/lang/meshfilesync/"
	hc.FileRoot = "data/meshfilesync/share/"
	hc.TempRoot = "data/meshfilesync/temp/"
	hc.MaxFileSize = 100.0
	hc.AutoAdd = false
	hc.AutoUpdate = false
	hc.AutoDelete = false
}

func (hc ModuleConfig) SaveConfig() error {

	mut.RLock()
	defer mut.RUnlock()

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

	mut.Lock()
	defer mut.Unlock()

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

func (hc *ModuleConfig) SetMaxFileSize(size float64) {

	mut.Lock()
	defer mut.Unlock()

	hc.MaxFileSize = size
}

func (hc *ModuleConfig) GetMaxFileSize() float64 {

	mut.RLock()
	defer mut.RUnlock()

	return hc.MaxFileSize
}

func (hc *ModuleConfig) SetAutoAdd(enable bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.AutoAdd = enable
}

func (hc *ModuleConfig) GetAutoAdd() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.AutoAdd
}

func (hc *ModuleConfig) SetAutoUpdate(enable bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.AutoUpdate = enable
}

func (hc *ModuleConfig) GetAutoUpdate() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.AutoUpdate
}

func (hc *ModuleConfig) SetAutoDelete(enable bool) {

	mut.Lock()
	defer mut.Unlock()

	hc.AutoDelete = enable
}

func (hc *ModuleConfig) GetAutoDelete() bool {

	mut.RLock()
	defer mut.RUnlock()

	return hc.AutoDelete
}
