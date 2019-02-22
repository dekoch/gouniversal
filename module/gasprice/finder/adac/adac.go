package adac

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/price"
	"github.com/dekoch/gouniversal/module/gasprice/station"
	"github.com/dekoch/gouniversal/shared/console"
)

func FindOnPage(st station.Station, raw string, prices *[]price.Price) error {

	var (
		err  error
		page string
		date time.Time
		pr   float64
	)

	p := make([]price.Price, 1)
	p[0].Date = time.Now()
	p[0].Station = st.UUID
	p[0].Source = "ADAC"
	p[0].Currency = "EUR"

	for i := 0; i <= 9; i++ {

		switch i {
		case 0:
			page = strings.Replace(raw, " ", "", -1)
			page = strings.Replace(page, "\t", "", -1)
			page = strings.Replace(page, "\n", "", -1)
			page = strings.Replace(page, "\r", "", -1)
			//fmt.Println(page)

		case 1:
			date, err = extractDate(page, "<tdclass=\"tal\">Diesel</td>")

		case 2:
			pr, err = extractPrice(page, "Diesel</td><tdclass=\"tal\">")

		case 3:
			p[0].AcquireDate = date
			p[0].Type = "Diesel"
			p[0].Price = pr
			*prices = append(*prices, p...)

		case 4:
			date, err = extractDate(page, "<tdclass=\"tal\">Super</td>")

		case 5:
			pr, err = extractPrice(page, "Super</td><tdclass=\"tal\">")

		case 6:
			p[0].AcquireDate = date
			p[0].Type = "Super"
			p[0].Price = pr
			*prices = append(*prices, p...)

		case 7:
			date, err = extractDate(page, "<tdclass=\"tal\">SuperE10</td>")

		case 8:
			pr, err = extractPrice(page, "SuperE10</td><tdclass=\"tal\">")

		case 9:
			p[0].AcquireDate = date
			p[0].Type = "Super E10"
			p[0].Price = pr
			*prices = append(*prices, p...)
		}

		if err != nil {
			console.Log(err, "")
			return err
		}
	}

	return nil
}

func extractDate(page string, entrypoint string) (time.Time, error) {

	var ret time.Time
	err := errors.New("date not found")

	lines := strings.SplitAfter(page, entrypoint)

	if len(lines) < 2 {
		return ret, err
	}

	tds := strings.Split(lines[1], "</td>")

	for _, td := range tds {

		td = strings.Replace(td, "<tdclass=\"tal\">", "", -1)

		// 12.02.201919:31:17
		if len(td) == 18 {

			dateDD := td[0:2]
			dateMM := td[3:5]
			dateYY := td[8:10]

			timeHH := td[10:12]
			timeMM := td[13:15]
			timeSS := td[16:18]

			ret, err = time.Parse("020106 150405", dateDD+dateMM+dateYY+" "+timeHH+timeMM+timeSS)
			return ret, err
		}
	}

	return ret, err
}

func extractPrice(page string, entrypoint string) (float64, error) {

	var ret float64
	err := errors.New("price not found")

	lines := strings.SplitAfter(page, entrypoint)

	if len(lines) < 2 {
		return ret, err
	}

	images := strings.Split(lines[1], "<imgsrc=\"/_common/img/info-test-rat/tanken/")

	price := ""

	func() {
		for _, image := range images {

			if strings.Contains(image, "</td>") {
				return
			}

			if strings.HasPrefix(image, "0.gif") {
				price += "0"
			} else if strings.HasPrefix(image, "1.gif") {
				price += "1"
			} else if strings.HasPrefix(image, "2.gif") {
				price += "2"
			} else if strings.HasPrefix(image, "3.gif") {
				price += "3"
			} else if strings.HasPrefix(image, "4.gif") {
				price += "4"
			} else if strings.HasPrefix(image, "5.gif") {
				price += "5"
			} else if strings.HasPrefix(image, "6.gif") {
				price += "6"
			} else if strings.HasPrefix(image, "7.gif") {
				price += "7"
			} else if strings.HasPrefix(image, "8.gif") {
				price += "8"
			} else if strings.HasPrefix(image, "9.gif") {
				price += "9"
			} else if strings.HasPrefix(image, "punkt.gif") {
				price += "."
			}
		}
	}()

	if price == "" {
		return 0.0, err
	}

	f, err := strconv.ParseFloat(price, 64)
	return f, err
}
