package globalOESPM

import (
	"gouniversal/modules/openespm/appConfig"
	"gouniversal/modules/openespm/deviceConfig"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/language"
)

const AppDataFolder = "data/openespm/app/"
const DeviceDataFolder = "data/openespm/device/"

var (
	UiConfig typesOESPM.UiConfig

	AppConfig    appConfig.AppConfig
	DeviceConfig deviceConfig.DeviceConfig

	Lang language.Language
)
