package gps

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/timeout"

	nmea "github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
)

type Gps struct {
	chanStart  chan bool
	chanFinish chan bool
	port       string
	baud       int
	sp         serial.Port
	nmea       []string
	pos        typenav.Pos
	watchdog   timeout.TimeOut
}

var (
	mut sync.RWMutex
)

func (gps *Gps) LoadConfig() {

	gps.chanStart = make(chan bool)
	gps.chanFinish = make(chan bool)

	go gps.readPort()
	go gps.job()
}

func (gps *Gps) OpenPort(port string, baud int, timeout int) error {

	mut.Lock()
	defer mut.Unlock()

	gps.port = port
	gps.baud = baud

	gps.watchdog.SetTimeOut(timeout)
	gps.watchdog.Enable(true)
	gps.watchdog.Reset()

	return gps.openPort()
}

func (gps *Gps) ClosePort() error {

	mut.Lock()
	defer mut.Unlock()

	return gps.closePort()
}

func (gps *Gps) openPort() error {

	fmt.Println("open port")

	var c serial.Config
	c.Name = gps.port
	c.Baud = gps.baud

	p, err := serial.OpenPort(&c)
	if err != nil {
		return err
	}

	gps.sp = *p

	return gps.sp.Flush()
}

func (gps *Gps) closePort() error {

	fmt.Println("close port")

	return gps.sp.Close()
}

func (gps *Gps) job() {

	intvl := time.Duration(20) * time.Millisecond
	tRead := time.NewTimer(intvl)

	for {

		select {
		case <-tRead.C:
			tRead.Stop()
			gps.chanStart <- true

		case <-gps.chanFinish:
			tRead.Reset(intvl)
		}
	}
}

func (gps *Gps) readPort() {

	var last string

	for {
		<-gps.chanStart

		func() {

			buf := make([]byte, 128)
			n, err := gps.sp.Read(buf)
			if err != nil {
				gps.closePort()
				gps.openPort()
			} else {

				//fmt.Println(string(buf[:n]))

				in := last + string(buf[:n])

				lines := strings.Split(in, "$")

				for _, l := range lines {

					if functions.IsEmpty(l) {
						continue
					}

					l = strings.TrimSpace(l)

					if validGPGGA("$" + l) {

						gps.watchdog.Reset()

						last = ""
						gps.parse("$" + l)
					} else {

						last = l
					}
				}
			}
		}()

		gps.chanFinish <- true
	}
}

func validGPGGA(s string) bool {

	if strings.HasPrefix(s, "$GPGGA") == false {
		return false
	}

	if strings.Contains(s, "*") == false {
		return false
	}

	end := strings.Split(s, "*")
	if len(end) >= 2 {

		checksum := strings.TrimSpace(end[1])
		if len(checksum) != 2 {

			return false
		}
	}

	return true
}

func (gps *Gps) parse(sentence string) {

	//fmt.Println(sentence)

	s, err := nmea.Parse(sentence)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch s.(type) {
	case nmea.GPGGA:

		g, ok := s.(nmea.GPGGA)
		if ok == false {
			return
		}

		mut.Lock()
		defer mut.Unlock()

		t := time.Now()
		gps.pos.Time = time.Date(t.Year(), t.Month(), t.Day(), g.Time.Hour, g.Time.Minute, g.Time.Second, g.Time.Millisecond, time.UTC)
		gps.pos.Lat = g.Latitude
		gps.pos.Lon = g.Longitude
		gps.pos.Fix = g.FixQuality
		gps.pos.Sat = g.NumSatellites
		gps.pos.HDOP = g.HDOP
		gps.pos.Ele = g.Altitude
		gps.pos.Valid = true
	}
}

func (gps *Gps) ResetTimeOut() {

	gps.watchdog.Reset()
}

func (gps *Gps) GetTimeOut() bool {

	return gps.watchdog.Elapsed()
}

func (gps *Gps) GetPos() typenav.Pos {

	mut.RLock()
	defer mut.RUnlock()

	return gps.pos
}
