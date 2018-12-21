package temphumv1x0

import (
	"encoding/json"

	"github.com/dekoch/gouniversal/module/openespm/respdevconfig"
	"github.com/dekoch/gouniversal/shared/console"
)

type DeviceConfig struct {
	Dev respdevconfig.RespDevConfig
}

func InitDeviceConfig() string {
	var c DeviceConfig
	c.Dev.Init()
	c.Dev.SetInterval(15.0*60.0, 1.0)

	b, err := json.Marshal(c)
	if err != nil {
		console.Log(err, "")
	}

	return string(b[:])
}

type AppConfig struct {
	DeviceUUID string
}

func InitAppConfig() string {
	var c AppConfig

	b, err := json.Marshal(c)
	if err != nil {
		console.Log(err, "")
	}

	return string(b[:])
}
