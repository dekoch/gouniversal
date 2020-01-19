package typemd

import (
	"github.com/dekoch/gouniversal/module/monmotion/lang"
)

type Page struct {
	Content string
	Lang    lang.LangFile
}

type Resolution struct {
	Width  int
	Height int
}
