package auftrag

import (
	"database/sql"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/google/uuid"
)

const TableName = "auftrag"

type Auftrag struct {
	ID        int
	UUID      string
	Name      string
	TypUUID   string
	IdentUUID string
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("id", sqlite3.TypeINTEGER, true, true)
	lyt.AddField("uuid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("name", sqlite3.TypeTEXT, false, false)
	lyt.AddField("typuuid", sqlite3.TypeTEXT, false, false)
	lyt.AddField("identuuid", sqlite3.TypeTEXT, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (ep *Auftrag) New(name string) string {

	u := uuid.Must(uuid.NewRandom())
	ep.UUID = u.String()

	ep.Name = name

	return ep.UUID
}

func (ep *Auftrag) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (uuid, name, typuuid, identuuid) values(?,?,?,?)", ep.UUID, ep.Name, ep.TypUUID, ep.IdentUUID)
	return err
}

func (ep *Auftrag) Load(name string, db *sql.DB) (bool, error) {

	var (
		err   error
		found bool
		id    string
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				rows, err = db.Query("SELECT id FROM `"+TableName+"` WHERE name=? ORDER BY id DESC LIMIT 0, 1", name)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&id)
				}

			case 3:
				found, err = ep.loadWithID(id, db)
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}

func (ep *Auftrag) loadWithID(id string, db *sql.DB) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = db.Query("SELECT id, uuid, name, typuuid, identuuid FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&ep.ID, &ep.UUID, &ep.Name, &ep.TypUUID, &ep.IdentUUID)

					found = true

					/*fmt.Println(ep.ID)
					fmt.Println(ep.UUID)
					fmt.Println(ep.Name)
					fmt.Println(ep.TypUUID)
					fmt.Println(ep.IdentUUID)*/
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}
