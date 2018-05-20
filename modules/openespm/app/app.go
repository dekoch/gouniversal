package app

import (
	"errors"
	"net/http"

	"github.com/dekoch/gouniversal/modules/openespm/app/SimpleSwitchV1x0/SimpleSwitchV1x0request"
	"github.com/dekoch/gouniversal/modules/openespm/app/SimpleSwitchV1x0/SimpleSwitchV1x0ui"
	"github.com/dekoch/gouniversal/modules/openespm/app/TempHumV1x0/TempHumV1x0request"
	"github.com/dekoch/gouniversal/modules/openespm/app/TempHumV1x0/TempHumV1x0ui"
	"github.com/dekoch/gouniversal/modules/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/navigation"
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
