package typesMesh

import (
	"github.com/dekoch/gouniversal/modules/mesh/network"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
)

type MeshError int

const (
	ErrNil MessageType = 1 + iota
	ErrNoConnection
	ErrServerDifferentMeshID
	ErrServerWrongMeshKey
	ErrServerWrongReceiver
	ErrServerEncryption
	ErrServerDecryption
	ErrClientDifferentMeshID
	ErrClientWrongMeshKey
	ErrClientWrongSender
	ErrClientEncryption
	ErrClientDecryption
)

type MessageType int

const (
	MessNil MessageType = 1 + iota
	MessAnnounce
	MessHello
	MessRAW
	MessMessenger
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
	Error    MessageType
}
