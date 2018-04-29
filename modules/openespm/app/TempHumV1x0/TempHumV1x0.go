package TempHumV1x0

import (
	"encoding/json"
	"fmt"
	"gouniversal/modules/openespm/respDevConfig"
)

type DeviceConfig struct {
	Dev respDevConfig.RespDevConfig
}

func InitDeviceConfig() string {
	var c DeviceConfig
	c.Dev.Init()
	c.Dev.SetInterval(15.0*60.0, 1.0)

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
