package serverList

import (
	"fmt"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
)

type ServerList struct {
	ServerList []serverInfo.ServerInfo
}

var (
	mut sync.RWMutex
)

func (sl *ServerList) Add(server serverInfo.ServerInfo) {

	mut.Lock()
	defer mut.Unlock()

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

	fmt.Println("add \"" + server.ID + "\"")

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

func (sl ServerList) Get() []serverInfo.ServerInfo {

	mut.RLock()
	defer mut.RUnlock()

	return sl.ServerList
}

func (sl *ServerList) Delete(id string) {

	mut.Lock()
	defer mut.Unlock()

	fmt.Println("delete \"" + id + "\"")

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

func (sl *ServerList) Clean(maxage float64) {

	mut.RLock()

	deleteID := ""

	for i := 0; i < len(sl.ServerList); i++ {

		if time.Since(sl.ServerList[i].TimeStamp).Hours() > maxage*24.0 {

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
