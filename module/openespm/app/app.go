package app

import (
	"errors"
	"net/http"

	"github.com/dekoch/gouniversal/module/openespm/app/simpleswitchv1x0/simpleswitchv1x0request"
	"github.com/dekoch/gouniversal/module/openespm/app/simpleswitchv1x0/simpleswitchv1x0ui"
	"github.com/dekoch/gouniversal/module/openespm/app/temphumv1x0/temphumv1x0request"
	"github.com/dekoch/gouniversal/module/openespm/app/temphumv1x0/temphumv1x0ui"
	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/navigation"
)

var UiAppList = [...]string{"SimpleSwitchV1x0", "TempHumV1x0"}
var DeviceAppList = [...]string{"SimpleSwitchV1x0", "TempHumV1x0"}

func Request(resp *typeoespm.Response, req *typeoespm.Request) {

	switch req.Device.App {
	case "SimpleSwitchV1x0":

		simpleswitchv1x0request.Request(resp, req)

	case "TempHumV1x0":

		temphumv1x0request.Request(resp, req)

	default:
		resp.Err = errors.New("app \"" + req.Device.App + "\" not found")
	}
}

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.IsNext("SimpleSwitchV1x0") {

		simpleswitchv1x0ui.Render(page, nav, r)

	} else if nav.IsNext("TempHumV1x0") {

		temphumv1x0ui.Render(page, nav, r)

	} else {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, "Render() "+nav.CurrentPath, nav.CurrentPath, nav.User.UUID)
	}
}
