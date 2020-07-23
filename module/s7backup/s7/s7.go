package s7

import (
	"time"

	"github.com/dekoch/gouniversal/module/s7backup/dbconfig"
	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/shared/io/s7conn"
)

func BackupPLC(comment string, pc *plcconfig.PLCConfig) error {

	var l []string

	for i := range pc.DB {
		l = append(l, pc.DB[i].UUID)
	}

	return BackupDB(l, comment, pc)
}

func BackupDB(uid []string, comment string, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			plc  s7conn.S7Conn
			conn *s7conn.Connection
		)

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				err = plc.AddPLC(pc.Address, pc.Rack, pc.Slot, 1, 1*time.Second, 1*time.Second)

			case 2:
				conn, err = plc.GetConnection(pc.Address)

			case 3:
				defer conn.Release()

			case 4:
				for i := range uid {

					err = backupDB(uid[i], comment, conn, pc)
					if err != nil {
						return
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func backupDB(uid, comment string, conn *s7conn.Connection, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			dc  dbconfig.DBConfig
			buf []byte
		)

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dc, err = pc.GetDB(uid)
				if dc.DBNo <= 0 ||
					dc.DBByteLength <= 0 {

					return
				}

			case 1:
				buf = make([]byte, dc.DBByteLength)
				err = conn.Client.AGReadDB(dc.DBNo, 0, dc.DBByteLength, buf)

			case 2:
				err = pc.SaveDB(global.Config.GetFileRoot(), uid, comment, buf)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func RestorePLC(pc *plcconfig.PLCConfig) error {

	var l []int

	for i := range pc.DB {

		dcs, err := pc.GetBackups(global.Config.GetFileRoot(), pc.DB[i].DBNo)
		if err != nil {
			return err
		}

		if len(dcs) > 0 {
			l = append(l, dcs[0].ID)
		}
	}

	return RestoreDB(l, pc)
}

func RestoreDB(id []int, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			plc  s7conn.S7Conn
			conn *s7conn.Connection
		)

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				err = plc.AddPLC(pc.Address, pc.Rack, pc.Slot, 1, 1*time.Second, 1*time.Second)

			case 1:
				conn, err = plc.GetConnection(pc.Address)

			case 2:
				defer conn.Release()

			case 3:
				for i := range id {

					err = restoreDB(id[i], conn, pc)
					if err != nil {
						return
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func restoreDB(id int, conn *s7conn.Connection, pc *plcconfig.PLCConfig) error {

	var err error

	func() {

		var (
			dc dbconfig.DBConfig
		)

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				dc, err = pc.LoadDB(global.Config.GetFileRoot(), id)

			case 1:
				err = conn.Client.AGWriteDB(dc.DBNo, 0, dc.DBByteLength, dc.DBData)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
