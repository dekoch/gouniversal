package global

import (
	"github.com/dekoch/gouniversal/module/openespm/appconfig"
	"github.com/dekoch/gouniversal/module/openespm/deviceconfig"
	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/shared/language"
)

const AppDataFolder = "data/openespm/app/"
const DeviceDataFolder = "data/openespm/device/"

var (
	UiConfig typeoespm.UiConfig

	AppConfig    appconfig.AppConfig
	DeviceConfig deviceconfig.DeviceConfig

	Lang language.Language
)
