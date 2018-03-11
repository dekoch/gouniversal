package types

import (
	"gouniversal/config"
	"sync"
)

type ProgramConfig struct {
	Header config.FileHeader
}

type User struct {
	UUID      string
	LoginName string
	Name      string
	PWDHash   string
	Groups    []string
	State     int
	Lang      string
	Comment   string
}

type UserConfigFile struct {
	Header config.FileHeader
	User   []User
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

type PageContent struct {
	Title string
}

type PageSettings struct {
	Title string
}

type Console struct {
	Mut   sync.Mutex
	Input string
}
