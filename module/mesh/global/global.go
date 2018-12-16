package global

import (
	"github.com/dekoch/gouniversal/module/mesh/keyfile"
	"github.com/dekoch/gouniversal/module/mesh/moduleConfig"
	"github.com/dekoch/gouniversal/module/mesh/networkConfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config        moduleConfig.ModuleConfig
	Lang          language.Language
	NetworkConfig networkConfig.NetworkConfig
	Keyfile       keyfile.Keyfile
)
