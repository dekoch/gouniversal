package typesMesh

import (
	"github.com/dekoch/gouniversal/modules/mesh/lang"
	"github.com/dekoch/gouniversal/modules/mesh/network"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
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
	Sender   serverInfo.ServerInfo
	Receiver serverInfo.ServerInfo
	Network  network.Network
	Message  ServerMessageContent
}

type Page struct {
	Content string
	Lang    lang.LangFile
}
