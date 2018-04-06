package typesOESPM

import (
	"encoding/json"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/shared/config"
	"net/url"
	"sync"
)

type RespType int

const (
	PLAIN RespType = 1 + iota
	JSON
	XML
)

type Response struct {
	Type    RespType
	Content string
	Status  int
	Err     error
}

type Device struct {
	UUID    string
	Key     string
	Name    string
	State   int
	Comment string
	App     string
	Config  string
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

type Request struct {
	UUID             string
	Key              string
	Values           url.Values
	Device           Device
	DeviceDataFolder string
}

type DeviceConfigFile struct {
	Header  config.FileHeader
	Devices []Device
}

type DeviceConfig struct {
	Mut  sync.Mutex
	File DeviceConfigFile
}

type JsonHeader struct {
	HeaderVersion float32
	AppName       string
	AppVersion    float32
}

type Page struct {
	Content string
	Lang    langOESPM.File
}

type UiConfig struct {
	Header      config.FileHeader
	AppFileRoot string
}
