package global

import (
	"gouniversal/modules/fileshare/moduleConfig"
	"gouniversal/shared/language"
	"gouniversal/shared/userToken"
)

var (
	Config moduleConfig.ModuleConfig
	Lang   language.Language
	Tokens userToken.UserToken
)
