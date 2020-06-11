package instauser

import (
	"database/sql"
	"time"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const TableName = "instauser"

type InstaUser struct {
	UserID    string
	UserName  string
	AccountID string
	Tag       string
	NewTags   []string
	Following bool
	Follow    time.Time
	Unfollow  time.Time
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("userid", sqlite3.TypeTEXT, true, false)
	lyt.AddField("username", sqlite3.TypeTEXT, false, false)
	lyt.AddField("accountid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("tag", sqlite3.TypeTEXT, false, false)
	lyt.AddField("following", sqlite3.TypeNUMERIC, false, false)
	lyt.AddField("follow", sqlite3.TypeDATE, false, false)
	lyt.AddField("unfollow", sqlite3.TypeDATE, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (inf *InstaUser) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (userid, username, accountid, tag, following, follow, unfollow) values(?,?,?,?,?,?,?)", inf.UserID, inf.UserName, inf.AccountID, inf.Tag, inf.Following, inf.Follow, inf.Unfollow)
	return err
}

func (inf *InstaUser) Load(userid string, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT userid, username, accountid, tag, following, follow, unfollow FROM `"+TableName+"` WHERE userid=?", userid)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&inf.UserID, &inf.UserName, &inf.AccountID, &inf.Tag, &inf.Following, &inf.Follow, &inf.Unfollow)
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

func Exists(userid string, dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT userid FROM `"+TableName+"` WHERE userid=?", userid)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				for dbconn.Rows.Next() {

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

func GetUsersFollowingBetween(fromdate, todate time.Time, dbconn *sqlite3.SQLite) ([]InstaUser, error) {

	var (
		err error
		ret []InstaUser
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				dbconn.Rows, err = dbconn.DB.Query("SELECT userid, username, accountid, tag, following, follow, unfollow FROM `"+TableName+"` WHERE following=TRUE AND follow BETWEEN ? AND ?", fromdate, todate)

			case 1:
				defer dbconn.Rows.Close()

			case 2:
				var n InstaUser

				for dbconn.Rows.Next() {

					err = dbconn.Rows.Scan(&n.UserID, &n.UserName, &n.AccountID, &n.Tag, &n.Following, &n.Follow, &n.Unfollow)
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
