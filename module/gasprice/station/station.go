package station

import (
	"errors"
	"sync"
)

type Station struct {
	UUID    string
	Name    string
	Company string
	Street  string
	City    string
	URL     string
}

type StationList struct {
	Stations []Station
}

var mut sync.RWMutex

func (sl *StationList) Add(st Station) {

	mut.Lock()
	defer mut.Unlock()

	for i, station := range sl.Stations {
		if station.UUID == st.UUID {

			sl.Stations[i] = st
			return
		}
	}

	sl.Stations = append(sl.Stations, st)
}

func (sl *StationList) GetList() []Station {

	mut.RLock()
	defer mut.RUnlock()

	return sl.Stations
}

func (sl *StationList) GetStation(uid string) (Station, error) {

	mut.RLock()
	defer mut.RUnlock()

	for _, station := range sl.Stations {

		if station.UUID == uid {
			return station, nil
		}
	}

	var station Station
	return station, errors.New("station not found")
}

func (sl *StationList) Remove(uid string) {

	mut.Lock()
	defer mut.Unlock()

	var l []Station

	for i := 0; i < len(sl.Stations); i++ {

		if uid != sl.Stations[i].UUID {
			l = append(l, sl.Stations[i])
		}
	}

	sl.Stations = l
}
