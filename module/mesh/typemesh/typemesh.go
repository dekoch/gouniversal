package typemesh

import (
	"github.com/dekoch/gouniversal/module/mesh/lang"
	"github.com/dekoch/gouniversal/module/mesh/network"
	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
)

type MessageType int

const (
	MessNil MessageType = 1 + iota
	MessAnnounce
	MessHello
	MessRAW
	MessMessenger
	MessFileSync
)

type ServerMessageContent struct {
	Type    MessageType
	Version float32
	Content []byte
}

type ServerMessage struct {
	Sender   serverinfo.ServerInfo
	Receiver serverinfo.ServerInfo
	Network  network.Network
	Message  ServerMessageContent
}

type Page struct {
	Content string
	Lang    lang.LangFile
}
