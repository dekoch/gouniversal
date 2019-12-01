package s7conn

import (
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/sbool"

	"github.com/dekoch/gos7"
)

const (
	// PLC Status
	s7CpuStatusUnknown = 0
	s7CpuStatusStop    = 4
	s7CpuStatusRun     = 8
)

type S7Conn struct {
	plcs []PLC
	mut  sync.Mutex
}

type PLCSettings struct {
	address              string
	rack, slot           int
	timeout, idleTimeout time.Duration
}

type PLC struct {
	settings    PLCSettings
	activeConn  int
	status      int
	connections []Connection
}

type Connection struct {
	connected bool
	no        int
	active    sbool.Sbool
	Handler   *gos7.TCPClientHandler
	Client    gos7.Client
}

func (pl *S7Conn) LoadConfig() {

}

func (pl *S7Conn) Exit() error {

	pl.mut.Lock()
	defer pl.mut.Unlock()

	for ip := range pl.plcs {

		for ic := range pl.plcs[ip].connections {

			err := pl.plcs[ip].connections[ic].close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pl *S7Conn) selectPLC(address string) *PLC {

	for i := range pl.plcs {

		if pl.plcs[i].settings.address == address {
			return &pl.plcs[i]
		}
	}

	var n PLC
	n.settings.address = address

	pl.plcs = append(pl.plcs, n)

	for i := range pl.plcs {

		if pl.plcs[i].settings.address == address {
			return &pl.plcs[i]
		}
	}

	return nil
}

func (pl *PLC) selectConnection() *Connection {

	for {
		for i := range pl.connections {

			if pl.connections[i].active.IsSet() == false {
				return &pl.connections[i]
			}
		}

		time.Sleep(1 * time.Millisecond)
	}
}

func (pl *S7Conn) AddPLC(address string, rack, slot, maxconn int, timeout, idleTimeout time.Duration) error {

	pl.mut.Lock()
	defer pl.mut.Unlock()

	plc := pl.selectPLC(address)
	plc.settings.rack = rack
	plc.settings.slot = slot
	plc.settings.timeout = timeout
	plc.settings.idleTimeout = idleTimeout

	for i := plc.activeConn; i < maxconn; i++ {

		var n Connection
		n.no = i
		plc.connections = append(plc.connections, n)

		plc.activeConn++
	}

	return nil
}

func (pl *S7Conn) GetConnection(address string) (*Connection, error) {

	pl.mut.Lock()
	defer pl.mut.Unlock()

	var (
		err    error
		status int
	)

	plc := pl.selectPLC(address)
	conn := plc.selectConnection()

	if conn.connected {

		status, err = conn.Client.PLCGetStatus()
		if err != nil {
			conn.connected = false
		}
	}

	if conn.connected == false {

		plc.status = 0

		err = conn.connect(plc.settings)
		if err != nil {
			return conn, err
		}

		status, err = conn.Client.PLCGetStatus()
		if err != nil {
			return conn, err
		}
	}

	if plc.status != status {

		plc.status = status

		printStatus(plc.settings.address, status)
	}

	conn.active.Set()

	return conn, nil
}

func (conn *Connection) Release() {

	conn.active.UnSet()
}

func (conn *Connection) connect(settings PLCSettings) error {

	if conn.connected {
		return nil
	}

	console.Output("("+settings.address+") connecting ("+strconv.Itoa(conn.no)+")", "s7conn")

	conn.Handler = gos7.NewTCPClientHandler(settings.address, settings.rack, settings.slot)
	conn.Handler.Timeout = settings.timeout
	conn.Handler.IdleTimeout = settings.idleTimeout

	err := conn.Handler.Connect()
	if err != nil {
		return err
	}

	conn.Client = gos7.NewClient(conn.Handler)

	conn.connected = true

	return nil
}

func (conn *Connection) close() error {

	return conn.Handler.Close()
}

func printStatus(address string, status int) {

	switch status {
	case s7CpuStatusUnknown:
		console.Output("("+address+") PLC UNKNOWN state", "s7conn")

	case s7CpuStatusStop:
		console.Output("("+address+") PLC STOP state", "s7conn")

	case s7CpuStatusRun:
		console.Output("("+address+") PLC RUN state", "s7conn")
	}
}
