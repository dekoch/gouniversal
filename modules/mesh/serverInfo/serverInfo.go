package serverInfo

import (
	"net"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/getPublicIP"
)

type ServerInfo struct {
	TimeStamp time.Time
	ID        string
	Port      int
	Address   []string
}

var (
	mut              sync.RWMutex
	pubAddrUpdInterv time.Duration
	timeUpdatedPA    time.Time
	publicAddress    string
)

func (si *ServerInfo) SetTimeStamp(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = t
}

func (si *ServerInfo) SetID(id string) {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	si.ID = id
}

func (si *ServerInfo) SetPort(port int) {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	si.Port = port
}

func (si *ServerInfo) GetPort() int {

	mut.RLock()
	defer mut.RUnlock()

	return si.Port
}

func (si *ServerInfo) SetPubAddrUpdInterv(interval int) {

	mut.Lock()
	defer mut.Unlock()

	pubAddrUpdInterv = time.Duration(interval)
}

func (si *ServerInfo) Update() {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	var empty []string
	si.Address = empty

	if pubAddrUpdInterv > 0 {
		si.publicAddress()
	}

	si.localAddresses()
}

func (si *ServerInfo) publicAddress() {

	if time.Since(timeUpdatedPA) > pubAddrUpdInterv*time.Minute {

		publicAddress = ""

		ip, err := getPublicIP.Get()
		if err == nil {
			publicAddress = ip
		}

		timeUpdatedPA = time.Now()
	}

	if publicAddress != "" {

		newAddress := make([]string, 1)
		newAddress[0] = publicAddress
		si.Address = append(si.Address, newAddress...)
	}
}

func (si *ServerInfo) localAddresses() {

	newAddress := make([]string, 1)

	// local addresses
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.IsLoopback() == false &&
					ipnet.IP.IsLinkLocalUnicast() == false &&
					ipnet.IP.IsLinkLocalMulticast() == false {

					newAddress[0] = ipnet.IP.String()

					//fmt.Println(newAddress[0])

					si.Address = append(si.Address, newAddress...)
				}
			}
		}
	}
}

func (si ServerInfo) Get() ServerInfo {

	mut.RLock()
	defer mut.RUnlock()

	return si
}
