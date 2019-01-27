package route

import (
	"bytes"
	"errors"
	"sync"

	"github.com/dekoch/gouniversal/module/gpsnav/geo"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"

	"github.com/dekoch/gouniversal/shared/io/file"
	gpx "github.com/twpayne/go-gpx"
)

type Route struct {
	filePath string
	wpt      []typenav.Pos
	nextWpt  int
}

var mut sync.RWMutex

func (rt *Route) Init(filePath string) {

	mut.Lock()
	defer mut.Unlock()

	rt.filePath = filePath
}

func (rt *Route) ReadFile() error {

	mut.Lock()
	defer mut.Unlock()

	b, err := file.ReadFile(rt.filePath)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	g, err := gpx.Read(r)
	if err != nil {
		return err
	}

	for _, p := range g.Wpt {

		var wpt typenav.Pos
		wpt.Time = p.Time
		wpt.Lat = p.Lat
		wpt.Lon = p.Lon
		wpt.Sat = int64(p.Sat)
		wpt.HDOP = p.HDOP
		wpt.Ele = p.Ele
		wpt.Name = p.Name
		wpt.Valid = true

		rt.wpt = append(rt.wpt, wpt)
	}

	return nil
}

func (rt *Route) StartRoute() {

	mut.Lock()
	defer mut.Unlock()

	rt.nextWpt = 0
}

func (rt *Route) SetWaypoint(no int) error {

	if no < 0 ||
		no >= len(rt.wpt) {
		return errors.New("invalid no")
	}

	mut.Lock()
	defer mut.Unlock()

	rt.nextWpt = no
	return nil
}

func (rt *Route) NextWaypoint() {

	mut.Lock()
	defer mut.Unlock()

	rt.nextWpt++
}

func (rt *Route) GetNextWaypoint() (typenav.Pos, error) {

	mut.RLock()
	defer mut.RUnlock()

	return rt.getWaypoint(rt.nextWpt)
}

func (rt *Route) GetWaypointNo() int {

	mut.RLock()
	defer mut.RUnlock()

	return rt.nextWpt
}

func (rt *Route) GetWaypoint(no int) (typenav.Pos, error) {

	mut.RLock()
	defer mut.RUnlock()

	return rt.getWaypoint(no)
}

func (rt *Route) getWaypoint(no int) (typenav.Pos, error) {

	if no < 0 ||
		no >= len(rt.wpt) {

		var ret typenav.Pos
		return ret, errors.New("invalid no")
	}

	return rt.wpt[no], nil
}

func (rt *Route) OnWaypoint(p typenav.Pos, distance float64) (bool, error) {

	mut.RLock()
	defer mut.RUnlock()

	wpt, err := rt.getWaypoint(rt.nextWpt)
	if err != nil {
		return false, err
	}

	var target geo.Geo
	target.SetTargetPos(wpt)
	target.SetCurrentPos(p)

	d, err := target.GetTargetDistance()
	if err != nil {
		return false, err
	}

	if d > distance {
		//fmt.Println(d)
		return false, nil
	}

	return true, nil
}
