package deviceConfig

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"
)

const configFilePath = "data/config/openespm/devices"

type Device struct {
	UUID       string
	RequestID  string
	RequestKey string
	Name       string
	State      int
	Comment    string
	App        string
	Config     string
}

func (d Device) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(d.Config), &v)
}

func (d *Device) Marshal(v interface{}) error {
	b, err := json.Marshal(v)
	if err == nil {
		d.Config = string(b[:])
	}

	return err
}

type DeviceConfigFile struct {
	Header  config.FileHeader
	Devices []Device
}

type DeviceConfig struct {
	Mut  sync.Mutex
	File DeviceConfigFile
}

func (dc DeviceConfig) SaveConfig() error {

	dc.File.Header = config.BuildHeader("devices", "devices", 1.0, "device config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newDevice := make([]Device, 1)

		newDevice[0].UUID = "test"
		newDevice[0].RequestID = "1234"
		newDevice[0].RequestKey = "1234"
		newDevice[0].Name = "Test"
		newDevice[0].State = 1 // active
		newDevice[0].App = "SimpleSwitchV1x0"

		dc.File.Devices = newDevice
	}

	b, err := json.Marshal(dc.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (dc *DeviceConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		dc.SaveConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &dc.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(dc.File.Header, "devices") == false {
		log.Fatal("wrong config")
	}

	return err
}
