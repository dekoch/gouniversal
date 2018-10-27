package serverList

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/shared/console"
)

type ServerList struct {
	ServerList []serverInfo.ServerInfo
	maxAge     float64
}

var (
	mut sync.RWMutex
)

func (sl *ServerList) SetMaxAge(maxage float64) {

	mut.Lock()
	defer mut.Unlock()

	sl.maxAge = maxage
}

func (sl *ServerList) Add(server serverInfo.ServerInfo) {

	mut.Lock()
	defer mut.Unlock()

	if sl.checkAge(server) == false {
		// server entry is too old
		return
	}

	for i := 0; i < len(sl.ServerList); i++ {
		// search server with ID
		if server.ID == sl.ServerList[i].ID {
			// check if newer
			if server.TimeStamp.After(sl.ServerList[i].TimeStamp) {
				// update
				sl.ServerList[i] = server
			}
			return
		}
	}

	console.Output("add \""+server.ID+"\"", "mesh")

	// add new server to list
	newServer := make([]serverInfo.ServerInfo, 1)
	newServer[0] = server
	sl.ServerList = append(sl.ServerList, newServer...)
}

func (sl *ServerList) AddList(serverList []serverInfo.ServerInfo) {

	for _, si := range serverList {

		sl.Add(si)
	}
}

func (sl *ServerList) Get() []serverInfo.ServerInfo {

	mut.RLock()
	defer mut.RUnlock()

	return sl.ServerList
}

func (sl *ServerList) GetWithID(id string) (serverInfo.ServerInfo, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(sl.ServerList); i++ {

		if id == sl.ServerList[i].ID {

			return sl.ServerList[i], nil
		}
	}

	var e serverInfo.ServerInfo
	return e, errors.New("server not found")
}

func (sl *ServerList) SetPrefAddress(server serverInfo.ServerInfo, addr string) {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(sl.ServerList); i++ {

		if server.ID == sl.ServerList[i].ID {

			sl.ServerList[i].SetPrefAddress(addr)
		}
	}
}

func (sl *ServerList) GetPrefAddress(server serverInfo.ServerInfo) string {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(sl.ServerList); i++ {

		if server.ID == sl.ServerList[i].ID {

			return sl.ServerList[i].GetPrefAddress()
		}
	}

	return ""
}

func (sl *ServerList) Delete(id string) {

	mut.Lock()
	defer mut.Unlock()

	console.Output("delete \""+id+"\"", "mesh")

	var l []serverInfo.ServerInfo
	n := make([]serverInfo.ServerInfo, 1)

	for i := 0; i < len(sl.ServerList); i++ {

		if id != sl.ServerList[i].ID {

			n[0] = sl.ServerList[i]

			l = append(l, n...)
		}
	}

	sl.ServerList = l
}

func (sl *ServerList) Clean() {

	mut.RLock()

	deleteID := ""

	for i := 0; i < len(sl.ServerList); i++ {

		if sl.checkAge(sl.ServerList[i]) == false {

			deleteID = sl.ServerList[i].ID
		}

		if deleteID != "" {

			mut.RUnlock()
			sl.Delete(deleteID)
			mut.RLock()

			deleteID = ""
		}
	}

	mut.RUnlock()
}

func (sl *ServerList) checkAge(server serverInfo.ServerInfo) bool {

	if time.Since(server.TimeStamp).Hours() > sl.maxAge*24.0 {
		return false
	}

	return true
}
