// package to load/save the module config

package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/google/uuid"
)

const configFilePath = "data/config/mesh/"

type ModuleConfig struct {
	Header           config.FileHeader
	UIFileRoot       string
	LangFileRoot     string
	PubAddrUpdInterv int // minutes (0=disabled)
	Server           serverinfo.ServerInfo
}

var (
	header config.FileHeader
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "mesh", ContentName: "mesh", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.Header = config.BuildHeaderWithStruct(header)

	hc.UIFileRoot = "data/ui/mesh/1.0/"
	hc.LangFileRoot = "data/lang/mesh/"

	hc.PubAddrUpdInterv = 30 // minutes

	// server defaults
	hc.Server.TimeStamp = time.Now()
	u := uuid.Must(uuid.NewRandom())
	hc.Server.ID = u.String() // UUID
	hc.Server.SetPort(9999)
}

func (hc ModuleConfig) SaveConfig() error {

	err := hc.CheckInput()
	if err != nil {
		return err
	}

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

	if hc.CheckInput() != nil {
		hc.loadDefaults()
	}

	hc.Server.SetPubAddrUpdInterv(hc.PubAddrUpdInterv)

	return err
}

func (hc *ModuleConfig) CheckInput() error {

	if functions.IsEmpty(hc.UIFileRoot) ||
		hc.PubAddrUpdInterv < 0 ||
		hc.PubAddrUpdInterv > 1440 ||
		functions.IsEmpty(hc.Server.ID) ||
		hc.Server.Port < 1 ||
		hc.Server.Port > 65535 {

		return errors.New("bad input")
	}

	return nil
}
