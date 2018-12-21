package global

import (
	"github.com/dekoch/gouniversal/module/mesh/keyfile"
	"github.com/dekoch/gouniversal/module/mesh/moduleconfig"
	"github.com/dekoch/gouniversal/module/mesh/networkconfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config        moduleconfig.ModuleConfig
	Lang          language.Language
	NetworkConfig networkconfig.NetworkConfig
	Keyfile       keyfile.Keyfile
)
