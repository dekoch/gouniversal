package parameter

import (
	"database/sql"
)

const TableName = "parameter"

type Parameter struct {
	ID            int
	UUID          string
	Typ           string
	Name          string
	Prozess       string
	ParameterNr   string
	ParameterUUID string
}

func (ep *Parameter) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (uuid, typ, name, prozess, parameternr, parameteruuid) values(?,?,?,?,?,?)", ep.UUID, ep.Typ, ep.Name, ep.Prozess, ep.ParameterNr, ep.ParameterUUID)
	return err
}

func (ep *Parameter) Load(uuid, prozess, parameternr string, db *sql.DB) (bool, error) {

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
				rows, err = db.Query("SELECT id FROM `"+TableName+"` WHERE uuid=? AND prozess=? AND parameternr=? ORDER BY id DESC LIMIT 0, 1", uuid, prozess, parameternr)

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

func (ep *Parameter) loadWithID(id string, db *sql.DB) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = db.Query("SELECT id, uuid, typ, name, prozess, parameternr, parameteruuid FROM `"+TableName+"` WHERE id=?", id)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&ep.ID, &ep.UUID, &ep.Typ, &ep.Name, &ep.Prozess, &ep.ParameterNr, &ep.ParameterUUID)

					found = true

					/*fmt.Println(ep.ID)
					fmt.Println(ep.UUID)
					fmt.Println(ep.Typ)
					fmt.Println(ep.Prozess)
					fmt.Println(ep.ParameterNr)
					fmt.Println(ep.ParameterUUID)*/
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}
