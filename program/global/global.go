package global

import (
	"gouniversal/program/groupConfig"
	"gouniversal/program/programTypes"
	"gouniversal/program/uiConfig"
	"gouniversal/program/userConfig"
	"gouniversal/shared/language"
)

type Global struct{}

var (
	Console programTypes.Console

	UiConfig uiConfig.UiConfig

	UserConfig userConfig.UserConfig

	GroupConfig groupConfig.GroupConfig

	Lang language.Language
)
