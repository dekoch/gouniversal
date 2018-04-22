package types

import "gouniversal/program/lang"

// Page is used to build a page and serve the selected language
type Page struct {
	Content string
	Lang    lang.File
}
