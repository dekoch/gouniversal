package global

import (
	"github.com/dekoch/gouniversal/module/picturex/moduleconfig"
	"github.com/dekoch/gouniversal/module/picturex/pairlist"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/token"
)

var (
	Config   moduleconfig.ModuleConfig
	Lang     language.Language
	Tokens   token.Token
	PairList pairlist.PairList
)
