package uiConfig

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/ui"

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

type UiConfigFile struct {
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

type UiConfig struct {
	Mut  sync.Mutex
	File UiConfigFile
}

func (c *UiConfig) SaveUiConfig() error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	c.File.Header = config.BuildHeader("ui", "ui", 1.0, "UI config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.File.UIEnabled = false
		c.File.Title = ""
		c.File.ProgramFileRoot = "data/ui/program/1.0/"
		c.File.StaticFileRoot = "data/ui/static/1.0/"

		c.File.HTTP.Enabled = true
		c.File.HTTP.Port = 8080

		c.File.HTTPS.Enabled = false
		c.File.HTTPS.Port = 443
		c.File.HTTPS.CertFile = "server.crt"
		c.File.HTTPS.KeyFile = "server.key"

		c.File.MaxGuests = 20
		c.File.MaxLoginAttempts = 10
		c.File.Recovery = false
	}

	b, err := json.Marshal(c.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (c *UiConfig) LoadConfig() {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.SaveUiConfig()
	}

	c.Mut.Lock()
	defer c.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &c.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(c.File.Header, "ui") == false {
		log.Fatal("wrong config \"" + configFilePath + "\"")
	}
}
