package typesFileshare

import (
	"net/url"

	"github.com/dekoch/gouniversal/modules/fileshare/lang"
)

type Page struct {
	Content string
	Lang    lang.LangFile
}

type Request struct {
	Values   url.Values
	ID       string
	Key      string
	FilePath string
}

type Response struct {
	Content  string
	FilePath string
	Status   int
	Err      error
}
