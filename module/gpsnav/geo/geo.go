package geo

import (
	"errors"
	"sync"

	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
)

type Geo struct {
	start   typenav.Pos
	target  typenav.Pos
	current typenav.Pos
}

var mut sync.RWMutex

func (geo *Geo) SetStartPos(p typenav.Pos) {

	mut.Lock()
	defer mut.Unlock()

	geo.start = p
}

func (geo *Geo) SetTargetPos(p typenav.Pos) {

	mut.Lock()
	defer mut.Unlock()

	geo.target = p
}

func (geo *Geo) SetCurrentPos(p typenav.Pos) {

	mut.Lock()
	defer mut.Unlock()

	geo.current = p
}

func (geo *Geo) IsStartPosValid() bool {

	mut.RLock()
	defer mut.RUnlock()

	return geo.start.Valid
}

func (geo *Geo) GetBearing() (float64, error) {

	mut.RLock()
	defer mut.RUnlock()

	if geo.start.Valid == false ||
		geo.current.Valid == false {
		return 0.0, errors.New("position not valid")
	}

	g := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	_, bearing := g.To(geo.start.Lat, geo.start.Lon, geo.current.Lat, geo.current.Lon)

	return bearing, nil
}

// GetTargetBearing returns the bearing relative from current bearing to target bearing
//   0.00 = on way
// -12.34 = target is left
// +12.34 = target is right
func (geo *Geo) GetTargetBearing() (float64, error) {

	mut.RLock()
	defer mut.RUnlock()

	if geo.start.Valid == false ||
		geo.target.Valid == false ||
		geo.current.Valid == false {
		return 0.0, errors.New("position not valid")
	}

	g := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	_, current := g.To(geo.start.Lat, geo.start.Lon, geo.current.Lat, geo.current.Lon)
	_, target := g.To(geo.start.Lat, geo.start.Lon, geo.target.Lat, geo.target.Lon)

	rel := current - target

	if rel > 180.0 {

		rel = rel - 360.0

	} else if rel < -180.0 {

		rel = rel + 360.0
	}

	return -rel, nil
}

func (geo *Geo) GetDistance() (float64, error) {

	mut.RLock()
	defer mut.RUnlock()

	if geo.start.Valid == false ||
		geo.current.Valid == false {
		return 0.0, errors.New("position not valid")
	}

	g := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	distance, _ := g.To(geo.start.Lat, geo.start.Lon, geo.current.Lat, geo.current.Lon)

	return distance, nil
}

func (geo *Geo) GetTargetDistance() (float64, error) {

	mut.RLock()
	defer mut.RUnlock()

	if geo.target.Valid == false ||
		geo.current.Valid == false {
		return 0.0, errors.New("position not valid")
	}

	g := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	distance, _ := g.To(geo.current.Lat, geo.current.Lon, geo.target.Lat, geo.target.Lon)

	return distance, nil
}
