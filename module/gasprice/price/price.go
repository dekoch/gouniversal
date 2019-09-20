package price

import (
	"database/sql"
	"time"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const TableName = "price"

type Price struct {
	Date        time.Time
	AcquireDate time.Time
	Station     string
	Type        string
	Price       float64
	Currency    string
	Source      string
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("id", sqlite3.TypeINTEGER, true, true)
	lyt.AddField("date", sqlite3.TypeDATE, false, false)
	lyt.AddField("acquiredate", sqlite3.TypeDATE, false, false)
	lyt.AddField("station", sqlite3.TypeTEXT, false, false)
	lyt.AddField("type", sqlite3.TypeTEXT, false, false)
	lyt.AddField("price", sqlite3.TypeREAL, false, false)
	lyt.AddField("currency", sqlite3.TypeTEXT, false, false)
	lyt.AddField("source", sqlite3.TypeTEXT, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (pr *Price) Save(tx *sql.Tx) error {

	_, err := tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (date, acquiredate, station, type, price, currency, source) values(?,?,?,?,?,?,?)", pr.Date, pr.AcquireDate, pr.Station, pr.Type, pr.Price, pr.Currency, pr.Source)
	return err
}

func LoadList(station, gastype string, fromdate, todate time.Time, db *sql.DB) ([]Price, error) {

	var (
		err error
		ret []Price
		pr  Price
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = db.Query("SELECT date, type, price FROM `"+TableName+"` WHERE station=? AND type=? AND date BETWEEN ? AND ?", station, gastype, fromdate, todate)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&pr.Date, &pr.Type, &pr.Price)
					if err != nil {
						return
					}

					ret = append(ret, pr)

					/*fmt.Println(pr.Date)
					fmt.Println(pr.Price)*/
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}
