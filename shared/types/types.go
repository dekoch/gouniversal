package types

import "gouniversal/program/lang"

// User stores all information about a single user
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

// Page is used to build a page and serve the selected language
type Page struct {
	Content string
	Lang    lang.File
}
