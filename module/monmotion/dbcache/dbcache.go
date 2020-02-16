package dbcache

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

//const DBFILE = "data/monmotion/cache.db3"

const DBFILE = ":memory:"

const TableName = "images"

var Cache DBCache

type DBCache struct {
	dbconn sqlite3.SQLite
	isOpen bool
}

var mut sync.RWMutex

func (ds *DBCache) Open() error {

	mut.Lock()
	defer mut.Unlock()

	if ds.isOpen {
		return nil
	}

	err := ds.dbconn.Open(DBFILE)
	if err != nil {
		return err
	}

	ds.isOpen = true

	return mdimg.LoadConfig(&ds.dbconn)
}

func (ds *DBCache) Close() error {

	mut.Lock()
	defer mut.Unlock()

	if ds.isOpen == false {
		return nil
	}

	return ds.dbconn.Close()
}

func (ds *DBCache) SaveImage(image *mdimg.MDImage) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 1:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 2:
				err = image.Save(ds.dbconn.Tx)

			case 3:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBCache) LoadImage(id string, image *mdimg.MDImage) error {

	mut.RLock()
	defer mut.RUnlock()

	found, err := image.Load(id, &ds.dbconn)
	if err != nil {
		return err
	}

	if found == false {
		return errors.New("id not found")
	}

	return nil
}

func (ds *DBCache) LoadImageInfo(id string, image *mdimg.MDImage) error {

	mut.RLock()
	defer mut.RUnlock()

	found, err := image.LoadInfo(id, &ds.dbconn)
	if err != nil {
		return err
	}

	if found == false {
		return errors.New("id not found")
	}

	return nil
}

func (ds *DBCache) DeleteImages(device string, state mdimg.ImageState, fromdate, todate time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 1:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 2:
				_, err = ds.dbconn.Tx.Exec("DELETE FROM `"+TableName+"` WHERE device=? AND state=? AND captured BETWEEN ? AND ?", device, state, fromdate, todate)

			case 3:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBCache) GetImageIDs(device string) ([]string, error) {

	mut.RLock()
	defer mut.RUnlock()

	var (
		err error
		ret []string
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE device=?", device)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				var id string

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&id)
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

func (ds *DBCache) GetImageIDsWithState(device string, state mdimg.ImageState) ([]string, error) {

	mut.RLock()
	defer mut.RUnlock()

	var (
		err error
		ret []string
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE device=? AND state=?", device, state)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				var id string

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&id)
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

func (ds *DBCache) GetImageIDsBetween(device string, fromdate, todate time.Time) ([]string, error) {

	mut.RLock()
	defer mut.RUnlock()

	var (
		err error
		ret []string
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE device=? AND captured BETWEEN ? AND ?", device, fromdate, todate)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				var id string

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&id)
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

func (ds *DBCache) SetStateToImages(device string, state mdimg.ImageState, fromdate, todate time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 1:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 2:
				_, err = ds.dbconn.Tx.Exec("UPDATE `"+TableName+"` SET state=? WHERE device=? AND captured BETWEEN ? AND ?", state, device, fromdate, todate)

			case 3:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
