package deviceconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
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
	} else {
		console.Log(err, "")
	}

	return err
}

type DeviceConfigFile struct {
	Header  config.FileHeader
	Devices []Device
}

type DeviceConfig struct {
	Mut  sync.RWMutex
	File DeviceConfigFile
}

func (c *DeviceConfig) SaveConfig() error {

	c.Mut.RLock()
	defer c.Mut.RUnlock()

	c.File.Header = config.BuildHeader("devices", "devices", 1.0, "device config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newDevice := make([]Device, 1)

		newDevice[0].UUID = "test"
		newDevice[0].RequestID = "1234"
		newDevice[0].RequestKey = "1234"
		newDevice[0].Name = "Test"
		newDevice[0].State = 1 // active
		newDevice[0].App = "SimpleSwitchV1x0"

		c.File.Devices = newDevice
	}

	b, err := json.Marshal(c.File)
	if err != nil {
		console.Log(err, "openESPM/deviceconfig.SaveConfig()")
	}

	err = file.WriteFile(configFilePath, b)

	return err
}

func (c *DeviceConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.SaveConfig()
	}

	c.Mut.Lock()
	defer c.Mut.Unlock()

	b, err := file.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "openESPM/deviceconfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &c.File)
	if err != nil {
		console.Log(err, "openESPM/deviceconfig.LoadConfig()")
	}

	if config.CheckHeader(c.File.Header, "devices") == false {
		console.Log("wrong config \""+configFilePath+"\"", "openESPM/deviceconfig.LoadConfig()")
	}

	return err
}

func (c *DeviceConfig) Add(d Device) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	newDevice := make([]Device, 1)

	newDevice[0] = d

	c.File.Devices = append(c.File.Devices, newDevice...)
}

func (c *DeviceConfig) Edit(uid string, d Device) error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.Devices); i++ {

		if uid == c.File.Devices[i].UUID {

			c.File.Devices[i] = d
			return nil
		}
	}

	return errors.New("Edit() device \"" + d.UUID + "\" not found")
}

func (c *DeviceConfig) Get(uid string) (Device, error) {

	c.Mut.RLock()
	defer c.Mut.RUnlock()

	for i := 0; i < len(c.File.Devices); i++ {

		if uid == c.File.Devices[i].UUID {

			return c.File.Devices[i], nil
		}
	}

	var d Device
	d.State = -1
	return d, errors.New("Get() device \"" + uid + "\" not found")
}

func (c *DeviceConfig) GetWithReqID(rid string) (Device, error) {

	c.Mut.RLock()
	defer c.Mut.RUnlock()

	for i := 0; i < len(c.File.Devices); i++ {

		if rid == c.File.Devices[i].RequestID {

			return c.File.Devices[i], nil
		}
	}

	var d Device
	d.State = -1
	return d, errors.New("GetWitheReqID() device \"" + rid + "\" not found")
}

func (c *DeviceConfig) List() []Device {

	c.Mut.RLock()
	defer c.Mut.RUnlock()

	return c.File.Devices
}

func (c *DeviceConfig) Delete(uid string) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	var dl []Device
	n := make([]Device, 1)

	for i := 0; i < len(c.File.Devices); i++ {

		if uid != c.File.Devices[i].UUID {

			n[0] = c.File.Devices[i]

			dl = append(dl, n...)
		}
	}

	c.File.Devices = dl
}
