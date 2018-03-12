package global

import (
	"gouniversal/program/lang"
	"gouniversal/program/programTypes"
)

type Global struct{}

var (
	Console programTypes.Console

	ProgramConfig programTypes.ProgramConfig

	UiConfig programTypes.UiConfig

	UserConfig programTypes.UserConfig

	GroupConfig programTypes.GroupConfig

	Lang lang.Global
)
