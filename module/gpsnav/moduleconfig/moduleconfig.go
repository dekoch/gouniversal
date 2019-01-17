package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/gpsnav/"

type gps struct {
	Port string
	Baud int
}

type gpx struct {
	FileRoot      string
	TrackDistance float64
}

type route struct {
	FilePath       string
	WptMaxDistance float64
}

type ModuleConfig struct {
	Header       config.FileHeader
	UIFileRoot   string
	LangFileRoot string
	Gps          gps
	Gpx          gpx
	Route        route
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "gpsnav", ContentName: "gpsnav", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/gpsnav/1.0/"
	hc.LangFileRoot = "data/lang/gpsnav/"

	hc.Gps.Port = "/dev/ttyACM0"
	hc.Gps.Baud = 19200

	hc.Gpx.FileRoot = "data/gpsnav/gpx/"
	hc.Gpx.TrackDistance = 1.0

	hc.Route.FilePath = "data/gpsnav/route.gpx"
	hc.Route.WptMaxDistance = 5.0
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

func (hc *ModuleConfig) SetGPSPort(port string) {

	mut.Lock()
	defer mut.Unlock()

	hc.Gps.Port = port
}

func (hc *ModuleConfig) GetGPSPort() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Gps.Port
}

func (hc *ModuleConfig) SetGPSBaud(baud int) {

	mut.Lock()
	defer mut.Unlock()

	hc.Gps.Baud = baud
}

func (hc *ModuleConfig) GetGPSBaud() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Gps.Baud
}

func (hc *ModuleConfig) GetGPXRoot() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Gpx.FileRoot
}

func (hc *ModuleConfig) GetGPXTrackDist() float64 {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Gpx.TrackDistance
}

func (hc *ModuleConfig) GetRouteFilePath() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Route.FilePath
}

func (hc *ModuleConfig) GetRouteWptMaxDist() float64 {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Route.WptMaxDistance
}
