package global

import (
	"gouniversal/program/types"
)

type Global struct{}

var (
	Console types.Console

	ProgramConfig types.ProgramConfig

	UiConfig types.UiConfig

	UserConfig types.UserConfig

	GroupConfig types.GroupConfig
)
