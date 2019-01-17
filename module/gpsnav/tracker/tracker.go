package tracker

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"github.com/dekoch/gouniversal/module/gpsnav/geo"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/io/file"
	gpx "github.com/twpayne/go-gpx"
)

type Tracker struct {
	filePath string
	distance float64
	enabled  bool
	autoSave bool
	wpt      []typenav.Pos
	geo.Geo
}

var mut sync.RWMutex

func (tr *Tracker) Init(filePath string, distance float64, autosave bool) {

	mut.Lock()
	defer mut.Unlock()

	tr.filePath = filePath
	tr.distance = distance
	tr.autoSave = autosave
}

func (tr *Tracker) Enable(b bool) {

	mut.Lock()
	defer mut.Unlock()

	tr.enabled = b
}

func (tr *Tracker) CheckPos(p typenav.Pos) error {

	mut.Lock()
	defer mut.Unlock()

	if tr.enabled == false {
		return nil
	}

	if tr.IsStartPosValid() == false {
		tr.SetStartPos(p)
	}

	tr.SetCurrentPos(p)

	d, err := tr.GetDistance()
	if err != nil {
		return err
	}

	if d > tr.distance {

		tr.SetStartPos(p)
		tr.newWaypoint(p)

		if tr.autoSave {
			return tr.writeFile()
		}
	}

	return nil
}

func (tr *Tracker) NewWaypoint(p typenav.Pos) {

	mut.Lock()
	defer mut.Unlock()

	tr.newWaypoint(p)
}

func (tr *Tracker) newWaypoint(p typenav.Pos) {

	tr.wpt = append(tr.wpt, p)
}

func (tr *Tracker) WriteFile() error {

	mut.RLock()
	defer mut.RUnlock()

	return tr.writeFile()
}

func (tr *Tracker) writeFile() error {

	var g gpx.GPX
	g.Version = "1.0"
	g.Creator = "ExpertGPS 1.1 - http://www.topografix.com"

	for i, p := range tr.wpt {

		var wpt gpx.WptType
		wpt.Time = p.Time
		wpt.Lat = p.Lat
		wpt.Lon = p.Lon

		switch p.Fix {
		case "1":
			wpt.Fix = "3d"

		case "2":
			wpt.Fix = "dpgs"

		case "3":
			wpt.Fix = "pps"

		default:
			wpt.Fix = "none"
		}

		wpt.Sat = int(p.Sat)
		wpt.HDOP = p.HDOP
		wpt.Ele = p.Ele

		if wpt.Name == "" {
			wpt.Name = strconv.Itoa(i)
		}

		wpt.Cmt = p.Comment

		g.Wpt = append(g.Wpt, &wpt)
	}

	var buf bytes.Buffer

	if err := g.WriteIndent(&buf, "", "  "); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(string(buf.Bytes()))

	return file.WriteFile(tr.filePath, buf.Bytes())
}
