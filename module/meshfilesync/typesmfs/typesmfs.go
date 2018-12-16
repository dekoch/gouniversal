package typesmfs

import (
	"github.com/dekoch/gouniversal/module/meshfilesync/lang"
	"github.com/dekoch/gouniversal/module/meshfilesync/syncfile"
)

type MessageType int

const (
	MessNil MessageType = 1 + iota
	MessList
	MessFileUploadReq
	MessFileUploadStart
	MessFileUpload
)

type Message struct {
	Type    MessageType
	Version float32
	Content []byte
	Error   error
}

type FileTransfer struct {
	FileInfo syncfile.SyncFile
	Content  []byte
}

type Page struct {
	Content string
	Lang    lang.LangFile
}
