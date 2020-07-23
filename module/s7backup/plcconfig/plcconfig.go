package plcconfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/s7backup/dbconfig"
	"github.com/dekoch/gouniversal/shared/io/fileinfo"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/google/uuid"
)

const TableName = "plc"

type PLCConfig struct {
	UUID         string
	Name         string
	Address      string
	Rack         int
	Slot         int
	MaxBackupCnt int
	DB           []dbconfig.DBConfig
	Created      time.Time
	Saved        time.Time
}

var mut sync.RWMutex

func NewPLC() (PLCConfig, error) {

	var ret PLCConfig
	ret.UUID = uuid.Must(uuid.NewRandom()).String()
	ret.Name = ret.UUID
	ret.Created = time.Now()
	ret.Address = "0.0.0.0"
	ret.Rack = 0
	ret.Slot = 2         // 2=CPU3xx, 1=CPU15xx
	ret.MaxBackupCnt = 0 // disabled
	return ret, nil
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("uuid", sqlite3.TypeTEXT, true, false)
	lyt.AddField("name", sqlite3.TypeTEXT, false, false)
	lyt.AddField("address", sqlite3.TypeTEXT, false, false)
	lyt.AddField("rack", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("slot", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("maxbackupcnt", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("created", sqlite3.TypeDATE, false, false)
	lyt.AddField("saved", sqlite3.TypeDATE, false, false)
	lyt.AddField("db", sqlite3.TypeTEXT, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (pc *PLCConfig) Save(path string) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	pc.Saved = time.Now()

	for i := range pc.DB {
		pc.DB[i].Saved = pc.Saved
	}

	func() {

		var (
			dbconn sqlite3.SQLite
			b      []byte
		)

		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				err = pc.cleanup()

			case 1:
				b, err = json.Marshal(pc.DB)

			case 2:
				err = dbconn.Open(path + pc.UUID + ".sqlite3")

			case 3:
				defer dbconn.Close()

			case 4:
				err = LoadConfig(&dbconn)

			case 5:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 6:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 7:
				_, err = dbconn.Tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (uuid, name, address, rack, slot, maxbackupcnt, created, saved, db) values(?,?,?,?,?,?,?,?,?)", pc.UUID, pc.Name, pc.Address, pc.Rack, pc.Slot, pc.MaxBackupCnt, pc.Created, pc.Saved, string(b))

			case 8:
				err = dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (pc *PLCConfig) Load(path, uid string) (bool, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err   error
		found bool
	)

	func() {

		var (
			dbconn sqlite3.SQLite
			db     string
		)

		for i := 0; i <= 7; i++ {

			switch i {
			case 0:
				_, err = os.Stat(path + uid + ".sqlite3")

			case 1:
				err = dbconn.Open(path + uid + ".sqlite3")

			case 2:
				defer dbconn.Close()

			case 3:
				dbconn.Rows, err = dbconn.DB.Query("SELECT uuid, name, address, rack, slot, maxbackupcnt, created, saved, db FROM `"+TableName+"` WHERE uuid=?", uid)

			case 4:
				defer dbconn.Rows.Close()

			case 5:
				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&pc.UUID, &pc.Name, &pc.Address, &pc.Rack, &pc.Slot, &pc.MaxBackupCnt, &pc.Created, &pc.Saved, &db)
					found = true
				}

			case 6:
				err = json.Unmarshal([]byte(db), &pc.DB)

			case 7:
				err = pc.cleanup()
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}

func (pc *PLCConfig) Delete(path string) error {

	mut.Lock()
	defer mut.Unlock()

	path += pc.UUID + ".sqlite3"

	if _, err := os.Stat(path); os.IsNotExist(err) == false {

		return os.Remove(path)
	}

	return nil
}

func (pc *PLCConfig) AddDB(dc dbconfig.DBConfig) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range pc.DB {

		if pc.DB[i].DBNo == dc.DBNo {
			pc.DB[i] = dc
			return nil
		}
	}

	pc.DB = append(pc.DB, dc)
	return nil
}

func (pc *PLCConfig) GetDB(uid string) (dbconfig.DBConfig, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := range pc.DB {

		if pc.DB[i].UUID == uid {
			return pc.DB[i], nil
		}
	}

	return dbconfig.DBConfig{}, errors.New("id not found")
}

func (pc *PLCConfig) DeleteDB(uid string) error {

	mut.Lock()
	defer mut.Unlock()

	return pc.deleteDB(uid)
}

func (pc *PLCConfig) deleteDB(uid string) error {

	var n []dbconfig.DBConfig

	for i := range pc.DB {

		if pc.DB[i].UUID != uid {
			n = append(n, pc.DB[i])
		}
	}

	pc.DB = n
	return nil
}

func (pc *PLCConfig) SaveDB(path, uid, comment string, b []byte) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var (
			dbconn sqlite3.SQLite
			ids    []int
		)

		for i := 0; i <= 9; i++ {

			switch i {
			case 0:
				err = dbconn.Open(path + pc.UUID + ".sqlite3")

			case 1:
				defer dbconn.Close()

			case 2:
				err = dbconfig.LoadConfig(&dbconn)

			case 3:
				if pc.MaxBackupCnt > 0 {

					for i := range pc.DB {

						if pc.DB[i].UUID == uid {

							ids, err = dbconfig.GetOldestIDs(pc.DB[i].DBNo, pc.MaxBackupCnt-1, &dbconn)
							if err != nil {
								return
							}
						}
					}
				}

			case 4:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 5:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 6:
				for i := range ids {
					err = dbconfig.Delete(ids[i], dbconn.Tx)
					if err != nil {
						return
					}
				}

			case 7:
				for i := range pc.DB {

					if pc.DB[i].UUID == uid {

						pc.DB[i].Comment = comment
						pc.DB[i].DBData = b
						err = pc.DB[i].SaveToDB(dbconn.Tx)
						if err != nil {
							return
						}
					}
				}

			case 8:
				err = dbconn.Tx.Commit()

			case 9:
				err = pc.cleanup()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (pc *PLCConfig) LoadDB(path string, id int) (dbconfig.DBConfig, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret dbconfig.DBConfig
	)

	func() {

		var dbconn sqlite3.SQLite

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				_, err = os.Stat(path + pc.UUID + ".sqlite3")

			case 1:
				err = dbconn.Open(path + pc.UUID + ".sqlite3")

			case 2:
				defer dbconn.Close()

			case 3:
				err = dbconfig.LoadConfig(&dbconn)

			case 4:
				_, err = ret.LoadFromDB(id, true, &dbconn)

			case 5:
				for i := range pc.DB {

					if pc.DB[i].UUID == ret.UUID {

						if pc.DB[i].DBByteLength != ret.DBByteLength ||
							pc.DB[i].DBByteLength != len(ret.DBData) {

							err = errors.New("byte length mismatch")
						}
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func (pc *PLCConfig) GetBackups(path string, dbno int) ([]dbconfig.DBConfig, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret []dbconfig.DBConfig
	)

	func() {

		var (
			dbconn sqlite3.SQLite
			ids    []int
		)

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				err = dbconn.Open(path + pc.UUID + ".sqlite3")

			case 1:
				defer dbconn.Close()

			case 2:
				err = dbconfig.LoadConfig(&dbconn)

			case 3:
				ids, err = dbconfig.GetLatestIDs(dbno, -1, &dbconn)

			case 4:
				for i := range ids {

					var n dbconfig.DBConfig
					_, err = n.LoadFromDB(ids[i], false, &dbconn)
					if err != nil {
						return
					}

					ret = append(ret, n)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func (pc *PLCConfig) cleanup() error {

	var dbsToDelete []string
	// search for dublicates
	for i := range pc.DB {

		cnt := 0

		for ii := range pc.DB {

			if pc.DB[i].DBNo == pc.DB[ii].DBNo {

				cnt++

				if cnt > 1 {
					dbsToDelete = append(dbsToDelete, pc.DB[ii].UUID)
				}
			}
		}
	}

	for i := range dbsToDelete {

		err := pc.deleteDB(dbsToDelete[i])
		if err != nil {
			return err
		}
	}

	for i := range pc.DB {

		pc.DB[i].DBData = []byte{}
	}

	return nil
}

func GetPLCList(path string) ([]string, error) {

	var ret []string

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return ret, err
		}
	}

	localFiles, err := fileinfo.Get(path, 0, false)
	if err != nil {
		return ret, err
	}

	for i := range localFiles {

		if strings.HasSuffix(localFiles[i].Name, ".sqlite3") == false {
			continue
		}

		localFiles[i].Name = strings.Replace(localFiles[i].Name, ".sqlite3", "", -1)

		ret = append(ret, localFiles[i].Name)
	}

	return ret, nil
}
