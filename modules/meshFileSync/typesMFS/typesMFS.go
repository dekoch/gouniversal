package typesMFS

import (
	"github.com/dekoch/gouniversal/modules/meshFileSync/lang"
	"github.com/dekoch/gouniversal/modules/meshFileSync/syncFile"
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
	FileInfo syncFile.SyncFile
	Content  []byte
}

type Page struct {
	Content string
	Lang    lang.LangFile
}
