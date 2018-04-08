package app

import (
	"errors"
	"gouniversal/modules/openespm/app/SimpleSwitch_v1_0/SimpleSwitch_v1_0_request"
	"gouniversal/modules/openespm/app/SimpleSwitch_v1_0/SimpleSwitch_v1_0_ui"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/navigation"
	"net/http"
)

const DataFolder = "data/config/openespm/"

func List() []string {

	s := []string{"SimpleSwitch_v1_0"}

	return s
}

func Request(resp *typesOESPM.Response, req *typesOESPM.Request) {

	req.DeviceDataFolder = DataFolder + req.UUID + "/"

	switch req.Device.App {
	case "SimpleSwitch_v1_0":

		SimpleSwitch_v1_0_request.Request(resp, req)

	default:
		resp.Err = errors.New("app \"" + req.Device.App + "\" not found")
	}
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.IsNext("SimpleSwitch_v1_0") {

		SimpleSwitch_v1_0_ui.Render(page, nav, r)
	}
}
