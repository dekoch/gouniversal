package openespm

import (
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/oespmGlobal"
	"gouniversal/modules/openespm/request"
	"gouniversal/modules/openespm/ui"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

	oespmGlobal.DeviceConfig.Mut.Lock()
	oespmGlobal.DeviceConfig.File = deviceManagement.LoadDevices()
	oespmGlobal.DeviceConfig.Mut.Unlock()

	request.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	request.Exit()
}
