package global

import (
	"gouniversal/program/groupConfig"
	"gouniversal/program/lang"
	"gouniversal/program/programTypes"
	"gouniversal/program/uiConfig"
	"gouniversal/program/userConfig"
)

type Global struct{}

var (
	Console programTypes.Console

	UiConfig uiConfig.UiConfig

	UserConfig userConfig.UserConfig

	GroupConfig groupConfig.GroupConfig

	Lang lang.Global
)
