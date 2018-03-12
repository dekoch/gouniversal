package programTypes

import (
	"gouniversal/shared/config"
	"gouniversal/shared/types"
	"sync"
)

type ProgramConfig struct {
	Header config.FileHeader
}

type UserConfigFile struct {
	Header config.FileHeader
	User   []types.User
}

type UserConfig struct {
	Mut  sync.Mutex
	File UserConfigFile
}

type Group struct {
	UUID         string
	Name         string
	State        int
	Comment      string
	CanModify    bool
	AllowedPages []string
}

type GroupConfigFile struct {
	Header config.FileHeader
	Group  []Group
}

type GroupConfig struct {
	Mut  sync.Mutex
	File GroupConfigFile
}

type UiConfig struct {
	Header   config.FileHeader
	FileRoot string
	Port     int
	Recovery bool
}

type Console struct {
	Mut   sync.Mutex
	Input string
}
