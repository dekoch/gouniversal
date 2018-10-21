package moduleConfig

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/modbustest/"

type Station struct {
	Active      bool
	ReadOffset  uint16
	WriteOffset uint16
}

type ModuleConfig struct {
	Header   config.FileHeader
	IP       string
	Port     string
	Station1 Station
	Station2 Station
}

var (
	header config.FileHeader
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "modbus", ContentName: "modbus", ContentVersion: 1.0, Comment: "modbus config file"}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.IP = "127.0.0.1"
	hc.Port = "502"

	hc.Station1.Active = true
	hc.Station1.ReadOffset = 0
	hc.Station1.WriteOffset = 0

	hc.Station2.Active = true
	hc.Station2.ReadOffset = 64
	hc.Station2.WriteOffset = 64
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
