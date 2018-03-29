package types

import "gouniversal/program/lang"

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

type Page struct {
	Content string
	Lang    lang.File
}
