package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/module/monmotion/core/coreconfig"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/monmotion/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	NoCores        int
	ExportFileRoot string
	Devices        []coreconfig.CoreConfig
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "monmotion", ContentName: "monmotion", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/monmotion/1.0/"
	hc.StaticFileRoot = "data/ui/monmotion/1.0/static/"
	hc.LangFileRoot = "data/lang/monmotion/"

	hc.ExportFileRoot = "data/monmotion/"
	hc.NoCores = 5
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

func (hc *ModuleConfig) SetExportFileRoot(path string) {

	mut.Lock()
	defer mut.Unlock()

	hc.ExportFileRoot = path
}

func (hc *ModuleConfig) GetExportFileRoot() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.ExportFileRoot
}

func (hc *ModuleConfig) SetNoCores(n int) {

	mut.Lock()
	defer mut.Unlock()

	hc.NoCores = n
}

func (hc *ModuleConfig) GetNoCores() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.NoCores
}

func (hc *ModuleConfig) GetDevices() []coreconfig.CoreConfig {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Devices
}

func (hc *ModuleConfig) AddNewDevice(dev coreconfig.CoreConfig) bool {

	mut.Lock()
	defer mut.Unlock()

	// check if device is already in list
	for i := range hc.Devices {

		if hc.Devices[i].Acquire.Device.Source == dev.Acquire.Device.Source {
			return false
		}
	}

	hc.Devices = append(hc.Devices, dev)

	return true
}

func (hc *ModuleConfig) SetEnabledDevices(uids []string) {

	mut.Lock()
	defer mut.Unlock()

	for i := range hc.Devices {

		hc.Devices[i].Enabled = false

		for _, selDev := range uids {

			if hc.Devices[i].UUID == selDev {

				hc.Devices[i].Enabled = true
			}
		}
	}
}

func (hc *ModuleConfig) GetEnabledDevices() []string {

	mut.RLock()
	defer mut.RUnlock()

	var ret []string

	for i := range hc.Devices {

		if hc.Devices[i].Enabled {
			ret = append(ret, hc.Devices[i].UUID)
		}
	}

	return ret
}

func (hc *ModuleConfig) GetDisabledDevices() []string {

	mut.RLock()
	defer mut.RUnlock()

	var ret []string

	for i := range hc.Devices {

		if hc.Devices[i].Enabled == false {
			ret = append(ret, hc.Devices[i].UUID)
		}
	}

	return ret
}

func (hc *ModuleConfig) GetDevice(uid string) (*coreconfig.CoreConfig, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := range hc.Devices {

		if hc.Devices[i].UUID == uid {

			return &hc.Devices[i], nil
		}
	}

	var n coreconfig.CoreConfig
	return &n, errors.New("device uuid not found")
}
