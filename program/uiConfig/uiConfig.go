package uiConfig

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"
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

func (uc UiConfig) SaveUiConfig() error {

	uc.File.Header = config.BuildHeader("ui", "ui", 1.0, "UI config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		uc.File.UIEnabled = false
		uc.File.Title = ""
		uc.File.ProgramFileRoot = "data/ui/program/1.0/"
		uc.File.StaticFileRoot = "data/ui/static/1.0/"

		uc.File.HTTP.Enabled = true
		uc.File.HTTP.Port = 8080

		uc.File.HTTPS.Enabled = false
		uc.File.HTTPS.Port = 443
		uc.File.HTTPS.CertFile = "server.crt"
		uc.File.HTTPS.KeyFile = "server.key"

		uc.File.MaxGuests = 20
		uc.File.MaxLoginAttempts = 10
		uc.File.Recovery = false
	}

	b, err := json.Marshal(uc.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (uc *UiConfig) LoadConfig() {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		uc.SaveUiConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &uc.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(uc.File.Header, "ui") == false {
		log.Fatal("wrong config \"" + configFilePath + "\"")
	}
}
