package global

import (
	"gouniversal/program/console"
	"gouniversal/program/groupConfig"
	"gouniversal/program/uiConfig"
	"gouniversal/program/userConfig"
	"gouniversal/shared/language"
)

type Global struct{}

var (
	Console console.Console

	UiConfig uiConfig.UiConfig

	UserConfig userConfig.UserConfig

	GroupConfig groupConfig.GroupConfig

	Lang language.Language
)
