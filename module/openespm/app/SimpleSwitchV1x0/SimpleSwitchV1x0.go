package SimpleSwitchV1x0

import (
	"encoding/json"

	"github.com/dekoch/gouniversal/shared/console"
)

type DeviceConfig struct {
	Switch bool
}

func InitDeviceConfig() string {
	var c DeviceConfig
	c.Switch = false

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
