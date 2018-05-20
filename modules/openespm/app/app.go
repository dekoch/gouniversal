package app

import (
	"errors"
	"gouniversal/modules/openespm/app/SimpleSwitchV1x0/SimpleSwitchV1x0request"
	"gouniversal/modules/openespm/app/SimpleSwitchV1x0/SimpleSwitchV1x0ui"
	"gouniversal/modules/openespm/app/TempHumV1x0/TempHumV1x0request"
	"gouniversal/modules/openespm/app/TempHumV1x0/TempHumV1x0ui"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/alert"
	"gouniversal/shared/navigation"
	"net/http"
)

var UiAppList = [...]string{"SimpleSwitchV1x0", "TempHumV1x0"}
var DeviceAppList = [...]string{"SimpleSwitchV1x0", "TempHumV1x0"}

func Request(resp *typesOESPM.Response, req *typesOESPM.Request) {

	switch req.Device.App {
	case "SimpleSwitchV1x0":

		SimpleSwitchV1x0request.Request(resp, req)

	case "TempHumV1x0":

		TempHumV1x0request.Request(resp, req)

	default:
		resp.Err = errors.New("app \"" + req.Device.App + "\" not found")
	}
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.IsNext("SimpleSwitchV1x0") {

		SimpleSwitchV1x0ui.Render(page, nav, r)

	} else if nav.IsNext("TempHumV1x0") {

		TempHumV1x0ui.Render(page, nav, r)

	} else {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, "Render() "+nav.CurrentPath, nav.CurrentPath, nav.User.UUID)
	}
}
