package simpleswitchv1x0request

// http://127.0.0.1:8080/request/?id=test&key=1234

import (
	"encoding/json"

	"github.com/dekoch/gouniversal/module/openespm/app/simpleswitchv1x0"
	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/shared/functions"
)

type appResp struct {
	Switch string
}

func Request(resp *typeoespm.Response, req *typeoespm.Request) {

	resp.Type = typeoespm.JSON

	// init new device
	if functions.IsEmpty(req.Device.Config) {
		req.Device.Config = simpleswitchv1x0.InitDeviceConfig()
	}

	// read device config
	var config simpleswitchv1x0.DeviceConfig
	err := req.Device.Unmarshal(&config)
	if err != nil {
		resp.Err = err
		return
	}

	// test
	//config.Switch = !config.Switch

	// build json response
	var js appResp
	if config.Switch {
		js.Switch = "on"
	} else {
		js.Switch = "off"
	}

	b, err := json.Marshal(js)
	if err != nil {
		resp.Err = err
	} else {
		resp.Content = string(b[:])
	}

	// write device config
	resp.Err = req.Device.Marshal(config)
}
