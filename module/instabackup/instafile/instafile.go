package instafile

import (
	"database/sql"
	"time"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const TableName = "instafile"

type InstaFile struct {
	UserID   string
	UserName string
	Added    time.Time
	FileID   string
	URL      string
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("id", sqlite3.TypeINTEGER, true, true)
	lyt.AddField("userid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("username", sqlite3.TypeTEXT, false, false)
	lyt.AddField("fileid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("added", sqlite3.TypeDATE, false, false)
	lyt.AddField("url", sqlite3.TypeTEXT, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (inf *InstaFile) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (userid, username, fileid, added, url) values(?,?,?,?,?)", inf.UserID, inf.UserName, inf.FileID, inf.Added, inf.URL)
	return err
}

func (inf *InstaFile) Load(id string, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = dbconn.DB.Query("SELECT userid, username, fileid, added, url FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&inf.UserID, &inf.UserName, &inf.FileID, &inf.Added, &inf.URL)
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

func (inf *InstaFile) Exists(dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE userid=? AND fileid=?", inf.UserID, inf.FileID)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

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
