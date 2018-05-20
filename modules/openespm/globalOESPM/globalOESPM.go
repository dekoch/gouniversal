package globalOESPM

import (
	"github.com/dekoch/gouniversal/modules/openespm/appConfig"
	"github.com/dekoch/gouniversal/modules/openespm/deviceConfig"
	"github.com/dekoch/gouniversal/modules/openespm/typesOESPM"
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
