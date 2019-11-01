package types

import "github.com/dekoch/gouniversal/program/lang"

// Page is used to build a page and serve the selected language
type Page struct {
	Content string
	Lang    lang.LangFile
}

type ExitMessage struct {
	Users []string
}
