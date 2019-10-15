package core

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/csv"
	"github.com/dekoch/gouniversal/module/gasprice/finder"
	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/price"
	"github.com/dekoch/gouniversal/module/gasprice/pricelist"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/fileinfo"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

var (
	chanCheckStart = make(chan bool)
)

func LoadConfig() {

	//splitLargeCSV()

	if global.Config.SaveToDB ||
		global.Config.LoadFromDB ||
		global.Config.ImportCSVtoDB {

		var dbconn sqlite3.SQLite

		err := dbconn.Open(global.Config.DBFile)
		if err != nil {
			console.Log(err, "")
			return
		}

		defer dbconn.Close()

		err = price.LoadConfig(&dbconn)
		if err != nil {
			console.Log(err, "")
		}

		if global.Config.ImportCSVtoDB {

			err = importCSVtoDB(global.Config.FileRoot, dbconn.DB, dbconn.Tx)
			if err != nil {
				console.Log(err, "")
				return
			}

			global.Config.ImportCSVtoDB = false
			err = global.Config.SaveConfig()
			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}

	go checkPrice()
	go job()

	if global.Config.UpdInterv == -1 {
		chanCheckStart <- true
	}
}

func Exit() {

}

func job() {

	intvl := global.Config.GetUpdInterval()
	timer := time.NewTimer(intvl)

	for {

		if intvl > 0 {

			select {
			case <-timer.C:
				chanCheckStart <- true
				timer.Reset(intvl)
			}
		} else {
			// wait until enabled
			time.Sleep(1 * time.Minute)
			intvl = global.Config.GetUpdInterval()
		}
	}
}

func checkPrice() {

	for {
		<-chanCheckStart

		fileName := time.Now().Format("2006-01-02")
		fileName += ".csv"

		for _, st := range global.Config.Stations.GetList() {

			if functions.IsEmpty(st.Name) ||
				functions.IsEmpty(st.URL) {

				continue
			}

			prices, err := finder.GetPrices(st)
			if err != nil {
				console.Log(err, "checkPrice()")
				continue
			}

			if global.Config.SaveToDB {

				err = savePricesToDB(prices)
				if err != nil {
					console.Log(err, "checkPrice()")
				}
			}

			if global.Config.SaveToCSV {

				for _, price := range prices {

					err = csv.Export(global.Config.FileRoot+fileName, price)
					if err != nil {
						console.Log(err, "checkPrice()")
					}
				}
			}
		}
	}
}

func savePricesToDB(prices []price.Price) error {

	var (
		err    error
		dbconn sqlite3.SQLite
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(global.Config.DBFile)

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
				// save
				for _, price := range prices {

					err = price.Save(dbconn.Tx)
					if err != nil {
						return
					}
				}

			case 5:
				err = dbconn.Tx.Commit()
			}
		}
	}()

	return err
}

func importCSVtoDB(fileroot string, db *sql.DB, tx *sql.Tx) error {

	files, err := fileinfo.Get(fileroot, 0, false)
	if err != nil {
		return err
	}

	for i := range files {

		if strings.HasSuffix(files[i].Name, ".csv") == false {
			continue
		}

		err = csvToDB(fileroot+files[i].Name, db, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func csvToDB(filepath string, db *sql.DB, tx *sql.Tx) error {

	var (
		err error
		pl  pricelist.PriceList
	)

	func() {

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				from := time.Now().AddDate(-999, 0, 0)
				pl, err = csv.Import(filepath, "*", "*", from)

			case 1:
				tx, err = db.Begin()

			case 2:
				defer func() {
					if err != nil {
						tx.Rollback()
					}
				}()

			case 3:
				// save
				for _, pr := range pl.Prices {

					err = pr.Save(tx)
					if err != nil {
						return
					}
				}

			case 4:
				err = tx.Commit()
			}
		}
	}()

	return err
}

func splitLargeCSV() {

	fileName := time.Now().Format("2006")
	fileName += ".csv"
	// split large file into multiple files
	if _, err := os.Stat(global.Config.FileRoot + fileName); os.IsNotExist(err) == false {

		err := csv.Split(global.Config.FileRoot+fileName, global.Config.FileRoot)
		if err != nil {
			fmt.Println(err)
		}

		err = os.Rename(global.Config.FileRoot+fileName, global.Config.FileRoot+fileName+".old")
		if err != nil {
			fmt.Println(err)
		}
	}
}
