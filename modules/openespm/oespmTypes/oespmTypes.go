package oespmTypes

import (
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
}

type Request struct {
	UUID   string
	Key    string
	Values url.Values
	Device Device
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
