package global

import (
	"github.com/dekoch/gouniversal/modules/fileshare/moduleConfig"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/userToken"
)

var (
	Config moduleConfig.ModuleConfig
	Lang   language.Language
	Tokens userToken.UserToken
)
