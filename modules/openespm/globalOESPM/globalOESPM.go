package globalOESPM

import (
	"gouniversal/modules/openespm/appConfig"
	"gouniversal/modules/openespm/deviceConfig"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
)

const AppDataFolder = "data/config/openespm/app/"
const DeviceDataFolder = "data/config/openespm/device/"

var (
	UiConfig typesOESPM.UiConfig

	AppConfig    appConfig.AppConfig
	DeviceConfig deviceConfig.DeviceConfig

	Lang langOESPM.Global
)
