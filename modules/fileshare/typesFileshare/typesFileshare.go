package typesFileshare

import (
	"gouniversal/modules/fileshare/lang"
	"net/url"
)

type Page struct {
	Content string
	Lang    lang.File
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
