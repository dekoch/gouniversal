package appManagement

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/app/SimpleSwitch_v1_0/request"
	"gouniversal/modules/openespm/typesOESPM"
)

const DataFolder = "data/config/openespm/"

func AppList() []string {

	s := []string{"SimpleSwitch_v1_0",
		"SimpleSwitch_v1_1"}

	return s
}

func Request(resp *typesOESPM.Response, req *typesOESPM.Request) {
	fmt.Println(req.UUID)
	fmt.Println(req.Key)
	fmt.Println(req.Device.Name)
	fmt.Println(req.Device.App)

	req.DeviceDataFolder = DataFolder + req.UUID + "/"

	switch req.Device.App {
	case "SimpleSwitch_v1_0":

		SimpleSwitch_v1_0_request.Request(resp, req)

	default:
		resp.Err = errors.New("app \"" + req.Device.App + "\" not found")
	}
}
