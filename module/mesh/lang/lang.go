package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type Server struct {
	Title            string
	Apply            string
	ID               string
	Port             string
	PubAddrUpdInterv string
	Addresses        string
}

type Network struct {
	Title            string
	Apply            string
	ID               string
	Settings         string
	AnnounceInterval string
	HelloInterval    string
	MaxClientAge     string
	AddServer        string
	Address          string
	Port             string
	Network          string
	LastSeen         string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header  config.FileHeader
	Title   string
	Server  Server
	Network Network
	Alert   Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangMesh", 1.0, "Language File")

	l.Title = "Mesh"

	l.Server.Title = "Server"
	l.Server.Apply = "Apply"
	l.Server.ID = "ID"
	l.Server.Port = "Port"
	l.Server.PubAddrUpdInterv = "Public Address Update Interval (0=disabled) [m]"
	l.Server.Addresses = "Addresses"

	l.Network.Title = "Network"
	l.Network.Apply = "Apply"
	l.Network.ID = "ID"
	l.Network.Settings = "Settings"
	l.Network.AnnounceInterval = "\"Announce\" Interval [s]"
	l.Network.HelloInterval = "\"Hello\" Interval [s]"
	l.Network.MaxClientAge = "Max Client Age [d]"
	l.Network.AddServer = "Add Server"
	l.Network.Address = "Address"
	l.Network.Port = "Port"
	l.Network.Network = "Network"
	l.Network.LastSeen = "Last Seen"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
