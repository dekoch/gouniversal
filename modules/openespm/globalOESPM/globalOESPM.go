package globalOESPM

import (
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
)

const AppDataFolder = "data/config/openespm/app/"
const DeviceDataFolder = "data/config/openespm/device/"

var (
	UiConfig typesOESPM.UiConfig

	AppConfig    typesOESPM.AppConfig
	DeviceConfig typesOESPM.DeviceConfig

	Lang langOESPM.Global
)
