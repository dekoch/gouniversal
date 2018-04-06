package SimpleSwitch_v1_0_request

// http://127.0.0.1:8080/request/?id=test&key=1234

import (
	"encoding/json"
	"gouniversal/modules/openespm/app/SimpleSwitch_v1_0"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/functions"
)

type appResp struct {
	Switch bool
}

func Request(resp *typesOESPM.Response, req *typesOESPM.Request) {

	resp.Type = typesOESPM.JSON

	// init new device
	if functions.IsEmpty(req.Device.Config) {
		req.Device.Config = SimpleSwitch_v1_0.InitConfig()
	}

	// read device config
	var config SimpleSwitch_v1_0.AppConfig
	err := req.Device.Unmarshal(&config)
	if err != nil {
		resp.Err = err
		return
	}

	// test
	//config.Switch = !config.Switch

	// build json response
	var js appResp
	js.Switch = config.Switch

	b, err := json.Marshal(js)
	if err != nil {
		resp.Err = err
	} else {
		resp.Content = string(b[:])
	}

	// write device config
	resp.Err = req.Device.Marshal(config)
}
