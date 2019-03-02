package csv

import (
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/price"
	"github.com/dekoch/gouniversal/shared/io/csv"
)

func Export(filepath string, pr price.Price) error {

	st, err := global.Config.Stations.GetStation(pr.Station)
	if err != nil {
		return err
	}

	row := make([]string, 13)
	row[0] = pr.Date.Format("2006-01-02")
	row[1] = pr.Date.Format("15:04:05")
	row[2] = pr.AcquireDate.Format("2006-01-02")
	row[3] = pr.AcquireDate.Format("15:04:05")
	row[4] = pr.Station
	row[5] = st.Name
	row[6] = st.Company
	row[7] = st.Street
	row[8] = st.City
	row[9] = pr.Source
	row[10] = pr.Type
	row[11] = strconv.FormatFloat(pr.Price, 'f', 3, 64)
	row[12] = pr.Currency

	return csv.AddRow(filepath, row)
}

func Import(filepath, uid, gastype string, from time.Time) (price.PriceList, error) {

	var (
		err error
		ret price.PriceList
		pr  price.Price
		t   time.Time
	)

	lines, err := csv.ReadAll(filepath)
	if err != nil {
		return ret, err
	}

	for _, line := range lines {

		err = nil

		func() {

			for i, val := range line {

				switch i {
				case 0:
					// date
					t, err = time.Parse("2006-01-02", val)
					if err == nil {
						pr.Date = t
					}

				case 1:
					// time
					t, err = time.Parse("15:04:05", val)
					if err == nil {
						pr.Date = pr.Date.Add(time.Duration(t.Hour())*time.Hour +
							time.Duration(t.Minute())*time.Minute +
							time.Duration(t.Second())*time.Second)
					}

				case 2:
					// acquire date
					t, err = time.Parse("2006-01-02", val)
					if err == nil {
						pr.AcquireDate = t
					}

				case 3:
					// acquire time
					t, err = time.Parse("15:04:05", val)
					if err == nil {
						pr.AcquireDate = pr.AcquireDate.Add(time.Duration(t.Hour())*time.Hour +
							time.Duration(t.Minute())*time.Minute +
							time.Duration(t.Second())*time.Second)
					}

				case 4:
					pr.Station = val

				case 5:
					// Name

				case 6:
					// Company

				case 7:
					// Street

				case 8:
					// City

				case 9:
					pr.Source = val

				case 10:
					pr.Type = val

				case 11:
					pr.Price, err = strconv.ParseFloat(val, 64)

				case 12:
					pr.Currency = val
				}

				if err != nil {
					return
				}
			}

			if pr.Date.Before(from) ||
				pr.Station != uid ||
				pr.Type != gastype {

				return
			}

			ret.Add(pr)
		}()
	}

	return ret, nil
}
