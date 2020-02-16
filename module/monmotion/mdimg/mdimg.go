package mdimg

import (
	"bytes"
	"database/sql"
	"errors"
	"image"
	"image/jpeg"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/dekoch/gouniversal/shared/mjpegavi1"
)

const TableName = "images"

type ImageState int

const (
	CACHE ImageState = 1
	SAVED ImageState = 2
)

type MDImage struct {
	ID          int
	Device      string
	State       ImageState
	Captured    time.Time
	Trigger     bool
	PreRecoding float64 // second
	Overrun     float64 // second
	Jpeg        []byte
	typemd.Resolution
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("id", sqlite3.TypeINTEGER, true, true)
	lyt.AddField("device", sqlite3.TypeTEXT, false, false)
	lyt.AddField("state", sqlite3.TypeINTEGER, false, false)
	lyt.AddField("captured", sqlite3.TypeDATE, false, false)
	lyt.AddField("trigger", sqlite3.TypeNUMERIC, false, false)
	lyt.AddField("prerecoding", sqlite3.TypeREAL, false, false)
	lyt.AddField("overrun", sqlite3.TypeREAL, false, false)
	lyt.AddField("selected", sqlite3.TypeNUMERIC, false, false)
	lyt.AddField("jpeg", sqlite3.TypeBLOB, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (md *MDImage) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (device, state, captured, trigger, prerecoding, overrun, jpeg) values(?,?,?,?,?,?,?)", md.Device, md.State, md.Captured, md.Trigger, md.PreRecoding, md.Overrun, md.Jpeg)
	return err
}

func (md *MDImage) Load(id string, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT device, state, captured, trigger, prerecoding, overrun, jpeg FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&md.Device, &md.State, &md.Captured, &md.Trigger, &md.PreRecoding, &md.Overrun, &md.Jpeg)
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

func (md *MDImage) LoadInfo(id string, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT device, state, captured, trigger, prerecoding, overrun FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&md.Device, &md.State, &md.Captured, &md.Trigger, &md.PreRecoding, &md.Overrun)
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

func (md *MDImage) EncodeImage(img image.Image) error {

	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return err
	}

	md.Jpeg = buf.Bytes()

	return nil
}

func (md *MDImage) GetImage() (image.Image, error) {

	b, err := mjpegavi1.Decode(md.Jpeg)
	if err != nil {
		var ret image.Image
		return ret, err
	}

	return jpeg.Decode(bytes.NewReader(b))
}

func GetCaptureTime(id string, dbconn *sqlite3.SQLite) (time.Time, error) {

	var (
		err error
		ret time.Time
	)

	func() {

		var found bool

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT captured FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&ret)
					found = true
				}

			case 3:
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
