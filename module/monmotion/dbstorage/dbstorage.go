package dbstorage

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/dbcache"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const DBFILE = "data/monmotion/storage.db3"
const TableName = "images"

var Stor DBStorage

type DBStorage struct {
	dbconn sqlite3.SQLite
	isOpen bool
}

type SequenceImage struct {
	ID       string
	Captured time.Time
}

type TriggerInfo struct {
	ID          string
	Device      string
	Captured    time.Time
	PreRecoding float64 // second
	Overrun     float64 // second
	Delete      bool
}

var mut sync.RWMutex

func (ds *DBStorage) Open() error {

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

func (ds *DBStorage) Close() error {

	mut.Lock()
	defer mut.Unlock()

	if ds.isOpen == false {
		return nil
	}

	return ds.dbconn.Close()
}

func (ds *DBStorage) LoadImage(id string, image *mdimg.MDImage) error {

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

func (ds *DBStorage) LoadImageInfo(id string, image *mdimg.MDImage) error {

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

func (ds *DBStorage) DeleteImages(device string, state mdimg.ImageState, fromdate, todate time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				ds.setFastOption()
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

func (ds *DBStorage) GetIDByTime(device string, captured time.Time) (string, error) {

	mut.RLock()
	defer mut.RUnlock()

	var (
		err error
		ret string
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE device=? AND captured=?", device, captured)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&ret)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func (ds *DBStorage) DeleteOldSequences(cnt int) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var tis []TriggerInfo

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				tis, err = ds.getTriggerInfos()

				if len(tis) <= cnt {
					return
				}

			case 1:
				err = ds.setFastOption()

			case 2:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 4:
				l := len(tis)

				if l > cnt {

					for n := 0; n < l-cnt; n++ {

						err = deleteSequence(tis, tis[n].ID, &ds.dbconn)
						if err != nil {
							return
						}
					}
				}

			case 5:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBStorage) DeleteSequence(triggerid string) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var tis []TriggerInfo

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				tis, err = ds.getTriggerInfos()

				if len(tis) == 0 {
					return
				}

			case 1:
				err = ds.setFastOption()

			case 2:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 4:
				for n := range tis {

					if tis[n].ID == triggerid {
						err = deleteSequence(tis, tis[n].ID, &ds.dbconn)
					}
				}

			case 5:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteSequence(trigger []TriggerInfo, id string, dbconn *sqlite3.SQLite) error {

	var err error

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				for n := range trigger {

					if trigger[n].ID == id {
						err = setSelected(trigger[n], true, dbconn)
					} else {
						err = setSelected(trigger[n], false, dbconn)
					}

					if err != nil {
						return
					}
				}

			case 1:
				// delete selected
				_, err = dbconn.Tx.Exec("DELETE FROM `" + TableName + "` WHERE selected=1")
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBStorage) getTriggerInfos() ([]TriggerInfo, error) {

	var (
		err error
		ret []TriggerInfo
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id, device, captured, prerecoding, overrun FROM `"+TableName+"` WHERE `trigger`=1 AND state=?", mdimg.SAVED)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				var ti TriggerInfo

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&ti.ID, &ti.Device, &ti.Captured, &ti.PreRecoding, &ti.Overrun)
					if err != nil {
						return
					}

					ret = append(ret, ti)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func setSelected(trigger TriggerInfo, selected bool, dbconn *sqlite3.SQLite) error {

	fromDate := trigger.Captured.Add(-time.Duration(trigger.PreRecoding) * time.Second)
	toDate := trigger.Captured.Add(time.Duration(trigger.Overrun) * time.Second)

	_, err := dbconn.Tx.Exec("UPDATE `"+TableName+"` SET selected=? WHERE device=? AND captured BETWEEN ? AND ?", selected, trigger.Device, fromDate, toDate)
	return err
}

func (ds *DBStorage) GetTriggerSI() ([]SequenceImage, error) {

	mut.RLock()
	defer mut.RUnlock()

	return ds.getTriggerSI()
}

func (ds *DBStorage) getTriggerSI() ([]SequenceImage, error) {

	var (
		err error
		ret []SequenceImage
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id, captured FROM `"+TableName+"` WHERE `trigger`=1 AND state=?", mdimg.SAVED)

			case 1:
				defer ds.dbconn.Rows.Close()

			case 2:
				var si SequenceImage

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&si.ID, &si.Captured)
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

func (ds *DBStorage) GetSequenceInfos(triggerid string, state mdimg.ImageState) ([]SequenceImage, error) {

	mut.RLock()
	defer mut.RUnlock()

	return ds.getSequenceInfos(triggerid, state)
}

func (ds *DBStorage) getSequenceInfos(triggerid string, state mdimg.ImageState) ([]SequenceImage, error) {

	var (
		err error
		ret []SequenceImage
	)

	func() {

		var triggerImg mdimg.MDImage

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				var found bool
				found, err = triggerImg.LoadInfo(triggerid, &ds.dbconn)

				if found == false {
					err = errors.New("id not found")
				}

			case 1:
				var fromDate, toDate time.Time
				fromDate = triggerImg.Captured.Add(-time.Duration(triggerImg.PreRecoding) * time.Second)
				toDate = triggerImg.Captured.Add(time.Duration(triggerImg.Overrun) * time.Second)

				ds.dbconn.Rows, err = ds.dbconn.DB.Query("SELECT id, captured FROM `"+TableName+"` WHERE device=? AND state=? AND captured BETWEEN ? AND ?", triggerImg.Device, state, fromDate, toDate)

			case 2:
				defer ds.dbconn.Rows.Close()

			case 3:
				var si SequenceImage

				for ds.dbconn.Rows.Next() {

					err = ds.dbconn.Rows.Scan(&si.ID, &si.Captured)
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

func (ds *DBStorage) SetStateToSequence(triggerid string, state mdimg.ImageState) error {

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var sis []SequenceImage

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				sis, err = ds.getSequenceInfos(triggerid, mdimg.CACHE)

				if len(sis) == 0 {
					return
				}

			case 1:
				err = ds.setFastOption()

			case 2:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 3:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 4:
				for _, si := range sis {

					_, err = ds.dbconn.Tx.Exec("UPDATE `"+TableName+"` SET state=? WHERE id=?", state, si.ID)
					if err != nil {
						return
					}
				}

			case 5:
				err = ds.dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBStorage) SaveBlock(device string, cacheids []string, rammaxage, dbmaxage time.Duration) error {

	if len(cacheids) == 0 {
		return nil
	}

	mut.Lock()
	defer mut.Unlock()

	var err error

	func() {

		var (
			img              mdimg.MDImage
			fromDate, toDate time.Time
		)

		for i := 0; i <= 9; i++ {

			switch i {
			case 0:
				err = dbcache.Cache.LoadImageInfo(cacheids[0], &img)
				if err != nil {
					return
				}

				fromDate = img.Captured

			case 1:
				err = dbcache.Cache.LoadImageInfo(cacheids[len(cacheids)-1], &img)
				if err != nil {
					return
				}

				toDate = img.Captured

			case 2:
				err = dbcache.Cache.SetStateToImages(device, mdimg.SAVED, fromDate, toDate)

			case 3:
				err = ds.setFastOption()

			case 4:
				ds.dbconn.Tx, err = ds.dbconn.DB.Begin()

			case 5:
				defer func() {
					if err != nil {
						ds.dbconn.Tx.Rollback()
					}
				}()

			case 6:
				_, err = ds.dbconn.Tx.Exec("DELETE FROM `"+TableName+"` WHERE device=? AND state=? AND captured BETWEEN ? AND ?", device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), toDate.Add(-dbmaxage))

			case 7:
				for _, id := range cacheids {

					err = dbcache.Cache.LoadImage(id, &img)
					if err != nil {
						return
					}

					img.State = mdimg.CACHE

					err = img.Save(ds.dbconn.Tx)
					if err != nil {
						return
					}
				}

			case 8:
				err = ds.dbconn.Tx.Commit()

			case 9:
				err = dbcache.Cache.DeleteImages(device, mdimg.SAVED, time.Now().AddDate(-999, 0, 0), toDate.Add(-rammaxage))
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ds *DBStorage) ExportSequence(triggerid string, path string) error {

	mut.RLock()
	defer mut.RUnlock()

	var err error

	func() {

		var (
			dir string
			sis []SequenceImage
		)

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				var t time.Time
				t, err = mdimg.GetCaptureTime(triggerid, &ds.dbconn)

				dir = t.Format("20060102_150405.0000") + "/"

			case 1:
				sis, err = ds.getSequenceInfos(triggerid, mdimg.SAVED)

			case 2:
				var (
					img   mdimg.MDImage
					found bool
				)

				for _, si := range sis {

					found, err = img.Load(si.ID, &ds.dbconn)
					if err != nil {
						return
					}

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

func (ds *DBStorage) Vacuum() error {

	mut.Lock()
	defer mut.Unlock()

	ds.setSafeOption()

	return ds.dbconn.Vacuum()
}

func (ds *DBStorage) setFastOption() error {

	ddl := `
		PRAGMA temp_store = "MEMORY";
		PRAGMA journal_mode = "MEMORY";
		PRAGMA secure_delete = "0";
		`

	_, err := ds.dbconn.DB.Exec(ddl)

	return err
}

func (ds *DBStorage) setSafeOption() error {

	ddl := `
		PRAGMA temp_store = "MEMORY";
		PRAGMA journal_mode = "DELETE";
		PRAGMA secure_delete = "0";
		`

	_, err := ds.dbconn.DB.Exec(ddl)

	return err
}
