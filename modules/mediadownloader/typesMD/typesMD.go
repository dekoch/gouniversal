package typesMD

import (
	"github.com/dekoch/gouniversal/modules/mediadownloader/lang"
)

type Page struct {
	Content string
	Lang    lang.LangFile
}

type DownloadFile struct {
	Url      string
	Filename string
}
