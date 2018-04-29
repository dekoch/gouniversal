package typesOESPM

import (
	"gouniversal/modules/openespm/appConfig"
	"gouniversal/modules/openespm/deviceConfig"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/shared/config"
	"net/url"
)

type RespType int

const (
	PLAIN RespType = 1 + iota
	JSON
	XML
)

type DefaultDevResp struct {
	Ver   float32
	Intvl float64
	Ds    bool
}

type Response struct {
	Type    RespType
	Content string
	Status  int
	Err     error
}

type Request struct {
	ID               string
	Key              string
	Values           url.Values
	Device           deviceConfig.Device
	DeviceDataFolder string
}

type JsonHeader struct {
	HeaderVersion float32
	AppName       string
	AppVersion    float32
}

type Page struct {
	Content string
	Lang    langOESPM.File
	App     appConfig.App
}

type UiConfig struct {
	Header      config.FileHeader
	AppFileRoot string
}
