package SimpleSwitchV1x0

import (
	"encoding/json"
	"fmt"
)

type DeviceConfig struct {
	Switch bool
}

func InitDeviceConfig() string {
	var c DeviceConfig
	c.Switch = false

	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	return string(b[:])
}
