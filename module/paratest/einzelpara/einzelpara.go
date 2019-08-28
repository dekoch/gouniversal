package einzelpara

import (
	"database/sql"
)

const TableName = "einzelpara"

type EinzelPara struct {
	ID   int
	UUID string
	Name string
	Wert string
}

func (ep *EinzelPara) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (uuid, name, wert) values(?,?,?)", ep.UUID, ep.Name, ep.Wert)
	return err
}

func (ep *EinzelPara) Load(uuid string, db *sql.DB) (bool, error) {

	var (
		err   error
		found bool
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = db.Query("SELECT id, uuid, name, wert FROM `"+TableName+"` WHERE uuid=? ORDER BY id DESC LIMIT 0, 1", uuid)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&ep.ID, &ep.UUID, &ep.Name, &ep.Wert)

					found = true

					/*fmt.Println(ep.ID)
					fmt.Println(ep.UUID)
					fmt.Println(ep.Name)
					fmt.Println(ep.Wert)*/
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}
