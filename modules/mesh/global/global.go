package global

import (
	"github.com/dekoch/gouniversal/modules/mesh/keyfile"
	"github.com/dekoch/gouniversal/modules/mesh/moduleConfig"
	"github.com/dekoch/gouniversal/modules/mesh/networkConfig"
)

var (
	Config        moduleConfig.ModuleConfig
	NetworkConfig networkConfig.NetworkConfig
	Keyfile       keyfile.Keyfile
)
