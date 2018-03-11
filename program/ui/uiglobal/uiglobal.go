package uiglobal

import (
	"gouniversal/program/lang"
)

var Lang lang.Global

type Page struct {
	Title   string
	Content string
	Lang    lang.File
}
