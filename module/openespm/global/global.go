package global

import (
	"github.com/dekoch/gouniversal/module/openespm/appConfig"
	"github.com/dekoch/gouniversal/module/openespm/deviceConfig"
	"github.com/dekoch/gouniversal/module/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/shared/language"
)

const AppDataFolder = "data/openespm/app/"
const DeviceDataFolder = "data/openespm/device/"

var (
	UiConfig typesOESPM.UiConfig

	AppConfig    appConfig.AppConfig
	DeviceConfig deviceConfig.DeviceConfig

	Lang language.Language
)
