package global

import (
	"github.com/dekoch/gouniversal/modules/mesh/keyfile"
	"github.com/dekoch/gouniversal/modules/mesh/moduleConfig"
	"github.com/dekoch/gouniversal/modules/mesh/networkConfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config        moduleConfig.ModuleConfig
	Lang          language.Language
	NetworkConfig networkConfig.NetworkConfig
	Keyfile       keyfile.Keyfile
)
