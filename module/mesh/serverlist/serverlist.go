package serverlist

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/shared/console"
)

type ServerList struct {
	ServerList []serverinfo.ServerInfo
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

func (sl *ServerList) Add(server serverinfo.ServerInfo) {

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
				if server.ExposePort == 0 {
					server.ExposePort = sl.ServerList[i].ExposePort
				}

				sl.ServerList[i] = server
			}
			return
		}
	}

	console.Output("add \""+server.ID+"\"", "mesh")

	// add new server to list
	sl.ServerList = append(sl.ServerList, server)
}

func (sl *ServerList) AddList(serverList []serverinfo.ServerInfo) {

	for _, si := range serverList {

		sl.Add(si)
	}
}

func (sl *ServerList) Get() []serverinfo.ServerInfo {

	mut.RLock()
	defer mut.RUnlock()

	return sl.ServerList
}

func (sl *ServerList) GetWithID(id string) (serverinfo.ServerInfo, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(sl.ServerList); i++ {

		if id == sl.ServerList[i].ID {

			return sl.ServerList[i], nil
		}
	}

	var e serverinfo.ServerInfo
	return e, errors.New("server not found")
}

func (sl *ServerList) SetPrefAddress(server serverinfo.ServerInfo, addr string) {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(sl.ServerList); i++ {

		if server.ID == sl.ServerList[i].ID {

			sl.ServerList[i].SetPrefAddress(addr)
		}
	}
}

func (sl *ServerList) GetPrefAddress(server serverinfo.ServerInfo) string {

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

	var l []serverinfo.ServerInfo

	for i := 0; i < len(sl.ServerList); i++ {

		if id != sl.ServerList[i].ID {
			l = append(l, sl.ServerList[i])
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

func (sl *ServerList) checkAge(server serverinfo.ServerInfo) bool {

	if time.Since(server.TimeStamp).Hours() > sl.maxAge*24.0 {
		return false
	}

	return true
}
