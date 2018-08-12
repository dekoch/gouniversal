package network

import (
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Network struct {
	TimeStamp        time.Time
	ID               string
	Port             int
	AnnounceInterval int
	HelloInterval    int
	MaxClientAge     float64
	Hash             string
}

var (
	mut      sync.RWMutex
	localKey []byte
)

func (sn *Network) SetKey(key []byte) {

	mut.Lock()
	defer mut.Unlock()

	bytes, err := bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err == nil {

		sn.Hash = string(bytes)
		localKey = key
	}
}

func (sn *Network) SetPort(port int) {

	mut.Lock()
	defer mut.Unlock()

	sn.Port = port
}

func (sn Network) GetPort() int {

	mut.RLock()
	defer mut.RUnlock()

	return sn.Port
}

func (sn Network) CheckID(id string) bool {

	mut.RLock()
	defer mut.RUnlock()

	if id == sn.ID {
		return true
	}

	return false
}

func (sn Network) CheckHashWithLocalKey(hash string) bool {

	mut.RLock()
	defer mut.RUnlock()

	err := bcrypt.CompareHashAndPassword([]byte(hash), localKey)
	return err == nil
}

func (sn Network) CheckKey(key []byte) bool {

	mut.RLock()
	defer mut.RUnlock()

	err := bcrypt.CompareHashAndPassword([]byte(sn.Hash), key)
	return err == nil
}

func (sn *Network) Update(newN Network) {

	mut.Lock()
	defer mut.Unlock()

	if newN.TimeStamp.After(sn.TimeStamp) {

		*sn = newN
	}
}

func (sn Network) Get() Network {

	mut.RLock()
	defer mut.RUnlock()

	return sn
}

func (sn Network) GetAnnounceInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(sn.AnnounceInterval) * time.Second
}

func (sn Network) GetHelloInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(sn.HelloInterval) * time.Second
}

func (sn Network) GetMaxClientAge() float64 {

	mut.RLock()
	defer mut.RUnlock()

	return sn.MaxClientAge
}
