package dbconfig

import (
	"database/sql"
	"time"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/google/uuid"
)

const TableName = "db"

type DBConfig struct {
	ID           int
	UUID         string
	DBNo         int
	DBByteLength int
	Name         string
	Created      time.Time
	Saved        time.Time
	Backup       time.Time
	DBData       []byte
}

func NewDB() (DBConfig, error) {

	var ret DBConfig
	ret.UUID = uuid.Must(uuid.NewRandom()).String()
	ret.Name = ret.UUID
	ret.Created = time.Now()
	return ret, nil
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("id", sqlite3.TypeINTEGER, true, true)
	lyt.AddField("uuid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("dbno", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("dbbytelenght", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("name", sqlite3.TypeTEXT, false, false)
	lyt.AddField("created", sqlite3.TypeDATE, false, false)
	lyt.AddField("saved", sqlite3.TypeDATE, false, false)
	lyt.AddField("backup", sqlite3.TypeDATE, false, false)
	lyt.AddField("dbdata", sqlite3.TypeBLOB, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func Delete(id int, tx *sql.Tx) error {

	_, err := tx.Exec("DELETE FROM `"+TableName+"` WHERE id=?", id)
	return err
}

func (dc *DBConfig) SaveToDB(tx *sql.Tx) error {

	dc.Backup = time.Now()

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (uuid, dbno, dbbytelenght, name, created, saved, backup, dbdata) values(?,?,?,?,?,?,?,?)", dc.UUID, dc.DBNo, dc.DBByteLength, dc.Name, dc.Created, dc.Saved, dc.Backup, dc.DBData)
	return err
}

func (dc *DBConfig) LoadFromDB(id int, withdata bool, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				if withdata {
					dbconn.Rows, err = dbconn.DB.Query("SELECT id, uuid, dbno, dbbytelenght, name, created, saved, backup, dbdata FROM `"+TableName+"` WHERE id=?", id)
				} else {
					dbconn.Rows, err = dbconn.DB.Query("SELECT id, uuid, dbno, dbbytelenght, name, created, saved, backup FROM `"+TableName+"` WHERE id=?", id)
				}

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					if withdata {
						err = dbconn.Rows.Scan(&dc.ID, &dc.UUID, &dc.DBNo, &dc.DBByteLength, &dc.Name, &dc.Created, &dc.Saved, &dc.Backup, &dc.DBData)
					} else {
						err = dbconn.Rows.Scan(&dc.ID, &dc.UUID, &dc.DBNo, &dc.DBByteLength, &dc.Name, &dc.Created, &dc.Saved, &dc.Backup)
					}

					found = true
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}

func GetLatestIDs(dbno, limit int, dbconn *sqlite3.SQLite) ([]int, error) {

	var (
		err error
		ret []int
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE dbno=? ORDER BY id DESC LIMIT 0, ?", dbno, limit)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					var id int

					err = dbconn.Rows.Scan(&id)
					if err != nil {
						return
					}

					ret = append(ret, id)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func GetOldestIDs(dbno, limit int, dbconn *sqlite3.SQLite) ([]int, error) {

	var (
		err error
		ret []int
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE dbno=? ORDER BY id DESC LIMIT ?, -1", dbno, limit)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					var id int

					err = dbconn.Rows.Scan(&id)
					if err != nil {
						return
					}

					ret = append(ret, id)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}
