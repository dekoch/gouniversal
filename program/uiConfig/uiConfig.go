package uiConfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/"

type UiHTTP struct {
	Enabled bool
	Port    int
}

type UiHTTPS struct {
	Enabled  bool
	Port     int
	CertFile string
	KeyFile  string
}

type UiConfig struct {
	Header           config.FileHeader
	UIEnabled        bool
	Title            string
	ProgramFileRoot  string
	StaticFileRoot   string
	HTTP             UiHTTP
	HTTPS            UiHTTPS
	MaxGuests        int
	MaxLoginAttempts int
	Recovery         bool
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "ui", ContentName: "ui", ContentVersion: 1.0, Comment: "UI config file"}
}

func (c *UiConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	c.UIEnabled = false
	c.Title = ""
	c.ProgramFileRoot = "data/ui/program/1.0/"
	c.StaticFileRoot = "data/ui/static/1.0/"

	c.HTTP.Enabled = true
	c.HTTP.Port = 8080

	c.HTTPS.Enabled = false
	c.HTTPS.Port = 443
	c.HTTPS.CertFile = "server.crt"
	c.HTTPS.KeyFile = "server.key"

	c.MaxGuests = 20
	c.MaxLoginAttempts = 10
	c.Recovery = false
}

func (c UiConfig) SaveConfig() error {

	mut.RLock()
	defer mut.RUnlock()

	c.Header = config.BuildHeaderWithStruct(header)

	b, err := json.Marshal(c)
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

func (c *UiConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		c.loadDefaults()
		c.SaveConfig()
	}

	mut.Lock()
	defer mut.Unlock()

	b, err := file.ReadFile(configFilePath + header.FileName)
	if err != nil {
		console.Log(err, "")
		c.loadDefaults()
	} else {
		err = json.Unmarshal(b, &c)
		if err != nil {
			console.Log(err, "")
			c.loadDefaults()
		}
	}

	if config.CheckHeader(c.Header, header.ContentName) == false {
		err = errors.New("wrong config \"" + configFilePath + header.FileName + "\"")
		console.Log(err, "")
		c.loadDefaults()
	}

	return err
}
