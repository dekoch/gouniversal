package typeoespm

import (
	"net/url"

	"github.com/dekoch/gouniversal/module/openespm/appconfig"
	"github.com/dekoch/gouniversal/module/openespm/deviceconfig"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/shared/config"
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
	Device           deviceconfig.Device
	DeviceDataFolder string
}

type JsonHeader struct {
	HeaderVersion float32
	AppName       string
	AppVersion    float32
}

type Page struct {
	Content string
	Lang    lang.LangFile
	App     appconfig.App
}

type UiConfig struct {
	Header      config.FileHeader
	AppFileRoot string
}
