package serverinfo

import (
	"net"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/getpublicip"
)

type ServerInfo struct {
	TimeStamp        time.Time
	ID               string
	Port             int
	ExposePort       int
	Address          []string
	publicAddress    string
	preferredAddress string
	manualAddress    string
	pubAddrUpdInterv time.Duration
	timeUpdatedPA    time.Time
}

var mut sync.RWMutex

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

func (si *ServerInfo) AddAddress(address string) {

	mut.Lock()
	defer mut.Unlock()

	for i := range si.Address {

		if si.Address[i] == address {
			return
		}
	}

	si.TimeStamp = time.Now()

	si.Address = append(si.Address, address)
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

func (si *ServerInfo) SetExposePort(port int) {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	si.ExposePort = port
}

func (si *ServerInfo) GetExposePort() int {

	mut.RLock()
	defer mut.RUnlock()

	return si.ExposePort
}

func (si *ServerInfo) SetPubAddrUpdInterv(interval int) {

	mut.Lock()
	defer mut.Unlock()

	si.pubAddrUpdInterv = time.Duration(interval)
}

func (si *ServerInfo) Update() {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	var empty []string
	si.Address = empty

	si.localAddresses()

	if si.pubAddrUpdInterv > 0 {
		si.pubAddress()
	}

	if si.publicAddress != "" {
		si.Address = append(si.Address, si.publicAddress)
	}

	if si.manualAddress != "" {
		si.Address = append(si.Address, si.manualAddress)
	}
}

func (si *ServerInfo) pubAddress() {

	if time.Since(si.timeUpdatedPA) > si.pubAddrUpdInterv*time.Minute {

		ip, err := getpublicip.Get()
		if err != nil {
			return
		}

		si.timeUpdatedPA = time.Now()

		si.publicAddress = ip
	}
}

func (si *ServerInfo) localAddresses() {

	// local addresses
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.IsLoopback() == false &&
					ipnet.IP.IsLinkLocalUnicast() == false &&
					ipnet.IP.IsLinkLocalMulticast() == false {

					si.Address = append(si.Address, ipnet.IP.String())
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

func (si *ServerInfo) SetPrefAddress(addr string) {

	mut.Lock()
	defer mut.Unlock()

	si.preferredAddress = addr
}

func (si *ServerInfo) GetPrefAddress() string {

	mut.RLock()
	defer mut.RUnlock()

	return si.preferredAddress
}

func (si *ServerInfo) SetManualAddress(addr string) {

	mut.Lock()
	defer mut.Unlock()

	si.TimeStamp = time.Now()

	si.manualAddress = addr
}

func (si *ServerInfo) GetManualAddress() string {

	mut.RLock()
	defer mut.RUnlock()

	return si.manualAddress
}
