package oespmTypes

import (
	"gouniversal/shared/config"
	"net/url"
	"sync"
)

type Request struct {
	UUID   string
	Key    string
	Values url.Values
}

type Device struct {
	UUID    string
	Key     string
	Name    string
	State   int
	Comment string
}

type DeviceConfigFile struct {
	Header  config.FileHeader
	Devices []Device
}

type DeviceConfig struct {
	Mut  sync.Mutex
	File DeviceConfigFile
}
