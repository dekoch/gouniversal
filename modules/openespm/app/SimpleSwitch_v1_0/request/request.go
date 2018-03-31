package SimpleSwitch_v1_0_request

// http://127.0.0.1:8080/request/?id=test&key=1234

import (
	"encoding/json"
	"gouniversal/modules/openespm/oespmTypes"
)

type appResp struct {
	Switch bool
}

func Request(resp *oespmTypes.Response, req *oespmTypes.Request) {

	resp.Type = oespmTypes.JSON

	var js appResp
	js.Switch = true

	b, err := json.Marshal(js)
	if err != nil {
		resp.Err = err
	} else {
		resp.Content = string(b[:])
	}
}
