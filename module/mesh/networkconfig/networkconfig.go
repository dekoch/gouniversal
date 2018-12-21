// package to load/save the network config from/to file

package networkconfig

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/network"
	"github.com/dekoch/gouniversal/module/mesh/serverlist"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/google/uuid"
)

const configFilePath = "data/config/mesh/"

type NetworkConfig struct {
	Header  config.FileHeader
	Network network.Network
	serverlist.ServerList
}

var (
	header config.FileHeader
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "network", ContentName: "network", ContentVersion: 1.0, Comment: ""}
}

func (hc *NetworkConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.Header = config.BuildHeaderWithStruct(header)

	hc.Network.TimeStamp = time.Now()
	u := uuid.Must(uuid.NewRandom())
	hc.Network.ID = u.String()       // UUID
	hc.Network.AnnounceInterval = 30 // seconds
	hc.Network.HelloInterval = 10    // seconds
	hc.Network.MaxClientAge = 30.0   // days
}

func (hc NetworkConfig) SaveConfig() error {

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

func (hc *NetworkConfig) LoadConfig() error {

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

	// check input
	if hc.Network.CheckConfig() == false {

		hc.loadDefaults()
	}

	return err
}
