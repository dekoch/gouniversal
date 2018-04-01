package SimpleSwitch_v1_0

import (
	"encoding/json"
	"fmt"
)

type AppConfig struct {
	Switch bool
}

func InitConfig() string {
	var ac AppConfig
	ac.Switch = false

	b, err := json.Marshal(ac)
	if err != nil {
		fmt.Println(err)
	}

	return string(b[:])
}
