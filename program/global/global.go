package global

import (
	"github.com/dekoch/gouniversal/program/console"
	"github.com/dekoch/gouniversal/program/groupConfig"
	"github.com/dekoch/gouniversal/program/uiConfig"
	"github.com/dekoch/gouniversal/program/userConfig"
	"github.com/dekoch/gouniversal/shared/language"
)

type Global struct{}

var (
	Console console.Console

	UiConfig uiConfig.UiConfig

	UserConfig userConfig.UserConfig

	GroupConfig groupConfig.GroupConfig

	Lang language.Language
)
