package modbustest

// only proof of concept

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/modbustest/moduleConfig"
	"github.com/dekoch/gouniversal/shared/console"

	"github.com/goburrow/modbus"
)

var (
	mConfig       moduleConfig.ModuleConfig
	mutConnection sync.Mutex
)

type modInput struct {
	Automatik         bool
	Einrichten        bool
	Stoerung          bool
	Start             bool
	Stop              bool
	Quitt             bool
	GSFahrt           bool
	Watchdog          bool
	ErgebnisLoeschen  bool
	StationAbgewaehlt bool
}

type modOutput struct {
	Stoerung     bool
	Bereit       bool
	Busy         bool
	Fertig       bool
	IO           bool
	NIO          bool
	Einrichten   bool
	WatchdogEcho bool
}

func bit(b byte, i int) bool {
	mask := []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
	value := b & mask[i]

	if value != 0 {
		return true
	}

	return false
}

func bitToValue(highbyte bool, bit int, val bool) uint16 {

	if val == true {

		mask := []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}

		value := mask[bit]

		out := []byte{0x00, 0x00}

		if highbyte {
			out[1] = value
		} else {
			out[0] = value
		}

		return binary.BigEndian.Uint16(out)
	}

	return 0
}

func readModbus(c moduleConfig.Station, client modbus.Client) (modInput, error) {

	fmt.Print("r ")

	var in modInput
	in.Automatik = false
	in.Einrichten = false
	in.Stoerung = false
	in.Start = false
	in.Stop = false
	in.Quitt = false
	in.GSFahrt = false
	in.Watchdog = false
	in.ErgebnisLoeschen = false
	in.StationAbgewaehlt = false

	// Read input register
	results, err := client.ReadInputRegisters(c.ReadOffset, 1)
	if err != nil {
		console.Log(err, "")
		return in, err
	}

	u := binary.BigEndian.Uint16(results)
	fmt.Print(u)

	in.Automatik = bit(results[1], 0)
	in.Einrichten = bit(results[1], 1)
	in.Stoerung = bit(results[1], 2)
	in.Start = bit(results[1], 3)
	in.Stop = bit(results[1], 4)
	in.Quitt = bit(results[1], 5)
	in.GSFahrt = bit(results[1], 6)
	in.Watchdog = bit(results[1], 7)
	in.ErgebnisLoeschen = bit(results[0], 0)
	in.StationAbgewaehlt = bit(results[0], 1)

	fmt.Println(in)

	return in, nil
}

func writeModbus(c moduleConfig.Station, o modOutput, client modbus.Client) error {

	fmt.Print("w ")

	var value uint16
	value = 0
	value += bitToValue(true, 0, o.Stoerung)
	value += bitToValue(true, 1, o.Bereit)
	value += bitToValue(true, 2, o.Busy)
	value += bitToValue(true, 3, o.Fertig)
	value += bitToValue(true, 4, o.IO)
	value += bitToValue(true, 5, o.NIO)
	value += bitToValue(true, 6, o.Einrichten)
	value += bitToValue(true, 7, o.WatchdogEcho)
	fmt.Print(value)

	fmt.Print(o)

	// Read input register
	results, err := client.WriteSingleRegister(c.WriteOffset, value)

	if err != nil {
		console.Log(err, "")
		return err
	}

	fmt.Println(results)

	return nil
}

func station(no int, c moduleConfig.Station, client modbus.Client) {

	fmt.Println(c.ReadOffset)
	fmt.Println(c.WriteOffset)

	var mIn modInput
	var mOut modOutput

	done := make(chan bool)

	step := 0
	oldStep := -1
	var err error

	for c.Active == true {

		mutConnection.Lock()

		mIn, err = readModbus(c, client)
		if err != nil {
			c.Active = false
		}

		mOut.WatchdogEcho = mIn.Watchdog

		if err == nil {

			if mIn.Quitt {
				mOut.Stoerung = false
			}

			if mIn.GSFahrt {
				mOut.Bereit = false
				mOut.Busy = false
				mOut.Fertig = false
				mOut.IO = false
				mOut.NIO = false
				mOut.Einrichten = false

				step = 0
			}

			switch step {
			case 0:
				mOut.Bereit = true
				mOut.Busy = false

				if mIn.ErgebnisLoeschen {
					mOut.Fertig = false
					mOut.IO = false
					mOut.NIO = false
				}

				if mIn.Start && mIn.ErgebnisLoeschen == false {
					step++
				}

			case 1:
				mOut.Bereit = false
				mOut.Busy = true

				mOut.Fertig = false
				mOut.IO = false
				mOut.NIO = false
				step++

			case 2:
				go func() {
					time.Sleep(time.Second * 5)

					done <- true
				}()

				step++

			case 3:
				go func() {
					if <-done {
						step++
					}
				}()

			case 4:
				mOut.IO = true
				step++

			case 5:
				mOut.Busy = false
				mOut.Fertig = true
				step++

			case 6:
				if mIn.Start == false {
					step = 0
				}

			default:
				step = 0
			}
		}

		err = writeModbus(c, mOut, client)
		if err != nil {
			c.Active = false
		}

		mutConnection.Unlock()

		if step != oldStep {
			oldStep = step

			fmt.Print("s")
			fmt.Print(no)
			fmt.Print(": step: ")
			fmt.Println(step)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func LoadConfig() {

	mConfig.LoadConfig()

	fmt.Println(mConfig.IP)
	fmt.Println(mConfig.Port)

	client := modbus.TCPClient(mConfig.IP + ":" + mConfig.Port)

	if mConfig.Station1.Active {
		go station(1, mConfig.Station1, client)
	}

	time.Sleep(1500 * time.Millisecond)

	if mConfig.Station2.Active {
		go station(2, mConfig.Station2, client)
	}
}
