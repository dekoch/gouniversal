package dbstorage

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const DBFILE = "data/monmotion/storage.db3"
const TableName = "images"

var mut sync.Mutex

type SequenceImage struct {
	ID       string
	Captured time.Time
}

func LoadConfig() error {

	mut.Lock()
	defer mut.Unlock()

	var dbconn sqlite3.SQLite

	err := dbconn.Open(DBFILE)
	if err != nil {
		return err
	}

	defer dbconn.Close()

	return mdimg.LoadConfig(&dbconn)
}

func SaveImages(images []mdimg.MDImage) error {

	mut.Lock()
	defer mut.Unlock()

	var (
		err    error
		dbconn sqlite3.SQLite
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 4:
				for _, img := range images {

					err = img.Save(dbconn.Tx)
					if err != nil {
						return
					}
				}

			case 5:
				err = dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func LoadImage(id string) (mdimg.MDImage, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret mdimg.MDImage
	)

	func() {

		var dbconn sqlite3.SQLite

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				var found bool
				found, err = ret.Load(id, &dbconn)

				if found == false {
					err = errors.New("id not found")
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func DeleteImages(device string, state mdimg.ImageState, fromdate, todate time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var dbconn sqlite3.SQLite

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 4:
				_, err = dbconn.Tx.Exec("DELETE FROM `"+TableName+"` WHERE device=? AND state=? AND captured BETWEEN ? AND ?", device, state, fromdate, todate)

			case 5:
				err = dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func SetStateToImages(device string, state mdimg.ImageState, fromdate, todate time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var dbconn sqlite3.SQLite

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 4:
				_, err = dbconn.Tx.Exec("UPDATE `"+TableName+"` SET state=? WHERE device=? AND captured BETWEEN ? AND ?", state, device, fromdate, todate)

			case 5:
				err = dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func GetTriggerIDs() ([]string, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret []string
	)

	func() {

		var (
			dbconn sqlite3.SQLite
			rows   *sql.Rows
		)

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				rows, err = dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE trigger=1 AND state=?", mdimg.SAVED)

			case 3:
				defer rows.Close()

			case 4:
				var id string

				for rows.Next() {

					err = rows.Scan(&id)
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

func GetSequenceInfos(id string) ([]SequenceImage, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret []SequenceImage
	)

	func() {

		var (
			dbconn     sqlite3.SQLite
			rows       *sql.Rows
			triggerImg mdimg.MDImage
		)

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				var found bool
				found, err = triggerImg.Load(id, &dbconn)

				if found == false {
					err = errors.New("id not found")
				}

			case 3:
				var fromDate, toDate time.Time
				fromDate = triggerImg.Captured.Add(-time.Duration(triggerImg.PreRecoding) * time.Second)
				toDate = triggerImg.Captured.Add(time.Duration(triggerImg.Overrun) * time.Second)

				/*fmt.Println(fromDate)
				fmt.Println(triggerImg.Captured)
				fmt.Println(toDate)*/

				rows, err = dbconn.DB.Query("SELECT id, captured FROM `"+TableName+"` WHERE device=? AND state=? AND captured BETWEEN ? AND ?", triggerImg.Device, mdimg.SAVED, fromDate, toDate)

			case 4:
				defer rows.Close()

			case 5:
				var si SequenceImage

				for rows.Next() {

					err = rows.Scan(&si.ID, &si.Captured)
					if err != nil {
						return
					}

					ret = append(ret, si)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func ExportSequence(triggerid string, ids []string, path string) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var (
			dbconn sqlite3.SQLite
			dir    string
		)

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				var t time.Time
				t, err = mdimg.GetCaptureTime(triggerid, &dbconn)

				dir = t.Format("20060102_150405.0000") + "/"

			case 3:
				var (
					img   mdimg.MDImage
					found bool
				)

				for _, id := range ids {

					found, err = img.Load(id, &dbconn)

					if found == false {
						continue
					}

					err = file.WriteFile(path+dir+img.Captured.Format("20060102_150405.0000")+".jpeg", img.Jpeg)
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

func Vacuum() error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var dbconn sqlite3.SQLite

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				err = dbconn.Open(DBFILE)

			case 1:
				defer dbconn.Close()

			case 2:
				err = dbconn.Vacuum()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
