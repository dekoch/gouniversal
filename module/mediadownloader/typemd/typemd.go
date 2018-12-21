package typemd

import (
	"github.com/dekoch/gouniversal/module/mediadownloader/lang"
)

type Page struct {
	Content string
	Lang    lang.LangFile
}

type DownloadFile struct {
	Url      string
	Filename string
}
