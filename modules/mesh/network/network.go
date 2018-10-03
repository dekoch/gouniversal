package network

import (
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"

	"golang.org/x/crypto/bcrypt"
)

type Network struct {
	TimeStamp        time.Time
	ID               string
	AnnounceInterval int     // seconds
	HelloInterval    int     // seconds (0=disabled)
	MaxClientAge     float64 // days
	Hash             string
}

var (
	mut      sync.RWMutex
	localKey []byte
)

func (sn *Network) SetTimeStamp(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	sn.TimeStamp = t
}

func (sn *Network) SetKey(key []byte) {

	mut.Lock()
	defer mut.Unlock()

	bytes, err := bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err == nil {

		sn.Hash = string(bytes)
		localKey = key
	}
}

func (sn *Network) GetID() string {

	mut.RLock()
	defer mut.RUnlock()

	return sn.ID
}

func (sn *Network) CheckID(id string) bool {

	mut.RLock()
	defer mut.RUnlock()

	if id == sn.ID {
		return true
	}

	return false
}

func (sn *Network) CheckHashWithLocalKey(hash string) bool {

	mut.RLock()
	defer mut.RUnlock()

	err := bcrypt.CompareHashAndPassword([]byte(hash), localKey)
	return err == nil
}

func (sn *Network) CheckKey(key []byte) bool {

	mut.RLock()
	defer mut.RUnlock()

	err := bcrypt.CompareHashAndPassword([]byte(sn.Hash), key)
	return err == nil
}

func (sn *Network) Update(net Network) {

	mut.Lock()
	defer mut.Unlock()

	if time.Now().After(net.TimeStamp) &&
		net.TimeStamp.After(sn.TimeStamp) &&
		net.CheckConfig() {

		console.Log("update network", "mesh")

		*sn = net
	}
}

func (sn Network) Get() Network {

	mut.RLock()
	defer mut.RUnlock()

	return sn
}

func (sn *Network) SetAnnounceInterval(sec int) {

	mut.Lock()
	defer mut.Unlock()

	sn.AnnounceInterval = sec
}

func (sn *Network) GetAnnounceInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(sn.AnnounceInterval) * time.Second
}

func (sn *Network) GetAnnounceIntervalInt() int {

	mut.RLock()
	defer mut.RUnlock()

	return sn.AnnounceInterval
}

func (sn *Network) SetHelloInterval(sec int) {

	mut.Lock()
	defer mut.Unlock()

	sn.HelloInterval = sec
}

func (sn *Network) GetHelloInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(sn.HelloInterval) * time.Second
}

func (sn *Network) GetHelloIntervalInt() int {

	mut.RLock()
	defer mut.RUnlock()

	return sn.HelloInterval
}

func (sn *Network) SetMaxClientAge(maxage float64) {

	mut.Lock()
	defer mut.Unlock()

	sn.MaxClientAge = maxage
}

func (sn *Network) GetMaxClientAge() float64 {

	mut.RLock()
	defer mut.RUnlock()

	return sn.MaxClientAge
}

func (sn *Network) CheckConfig() bool {

	if functions.IsEmpty(sn.ID) == false &&
		sn.AnnounceInterval >= 1 &&
		sn.AnnounceInterval <= 900 &&
		sn.HelloInterval >= 0 &&
		sn.HelloInterval <= 900 &&
		sn.MaxClientAge >= 1.0 &&
		sn.MaxClientAge <= 365.0 {

		return true
	}

	return false
}
