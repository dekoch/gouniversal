package price

import (
	"time"
)

type Price struct {
	Date        time.Time
	AcquireDate time.Time
	Station     string
	Type        string
	Price       float64
	Currency    string
	Source      string
}

type PriceList struct {
	Prices []Price
}

func (pl *PriceList) Add(pr Price) {

	pl.Prices = append(pl.Prices, pr)
}

func (pl *PriceList) AddList(prs []Price) {

	pl.Prices = append(pl.Prices, prs...)
}

func (pl *PriceList) GetStationUUIDs(gastype string) []string {

	var (
		ret     []string
		missing bool
	)

	for i := len(pl.Prices) - 1; i >= 0; i-- {

		if pl.Prices[i].Type != gastype {
			continue
		}

		missing = true

		for r := len(ret) - 1; r >= 0; r-- {

			if ret[r] == pl.Prices[i].Station {

				missing = false
			}
		}

		if missing {
			ret = append(ret, pl.Prices[i].Station)
		}
	}

	return ret
}

func (pl *PriceList) GetList() []Price {

	return pl.Prices
}

func (pl *PriceList) GetFromTimeSpan(uid, gastype string, from, to time.Time) []Price {

	var ret []Price

	func() {

		for i := len(pl.Prices) - 1; i >= 0; i-- {

			if pl.Prices[i].Station != uid ||
				pl.Prices[i].Type != gastype {

				continue
			}

			if pl.Prices[i].Date.After(from) &&
				pl.Prices[i].Date.Before(to) {

				ret = append(ret, pl.Prices[i])
			}
		}
	}()

	return ret
}
