package pagehome

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/csv"
	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/lang"
	"github.com/dekoch/gouniversal/module/gasprice/price"
	"github.com/dekoch/gouniversal/module/gasprice/pricelist"
	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/fileInfo"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/timeout"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:GasPrice:Home", page.Lang.Home.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang         lang.Home
		CurrentPrice template.HTML
		CmbStation   template.HTML
		CmbGasType   template.HTML
		Currency     template.JS
		Labels       template.JS
		Datasets     template.JS
	}
	var c Content

	c.Lang = page.Lang.Home

	id := ""
	selStation := r.FormValue("Station")
	stations := global.Config.Stations.GetList()

	sort.Slice(stations, func(i, j int) bool { return stations[i].City < stations[j].City })

	if selStation != "" {

		id = selStation

	} else if id == "" && len(stations) > 0 {

		id = stations[0].UUID
	}

	cmbStation := "<select name=\"Station\">"

	for _, station := range stations {

		cmbStation += "<option value=\"" + station.UUID + "\""

		if station.UUID == id {
			cmbStation += " selected"
		}

		cmbStation += ">" + station.City + " - " + station.Street + " - " + station.Name + "</option>"
	}
	cmbStation += "</select>"
	c.CmbStation = template.HTML(cmbStation)

	gasType := ""
	selGasType := r.FormValue("GasType")
	gasTypes := global.Config.GetGasTypes()

	if selGasType != "" {

		gasType = selGasType

	} else if gasType == "" && len(gasTypes) > 0 {

		gasType = gasTypes[0]
	}

	cmbGasType := "<select name=\"GasType\">"

	for _, gasType := range gasTypes {

		cmbGasType += "<option value=\"" + gasType + "\""

		if gasType == selGasType {
			cmbGasType += " selected"
		}

		cmbGasType += ">" + gasType + "</option>"
	}
	cmbGasType += "</select>"
	c.CmbGasType = template.HTML(cmbGasType)

	t := time.Now()
	from := t

	edit := r.FormValue("edit")

	switch edit {
	default:
		from = from.AddDate(0, 0, -1)

	case "7days":
		from = from.AddDate(0, 0, -7)

	case "30days":
		from = from.AddDate(0, 0, -30)

	case "alldays":
		from = from.AddDate(-999, 0, 0)
	}

	var (
		err    error
		prices []price.Price
	)

	if global.Config.LoadFromDB {
		prices, err = loadFromDB(id, gasType, from, t)
	}

	if global.Config.LoadFromCSV {
		prices, err = loadFromCSV(id, gasType, from, t)
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		return
	}

	labels := ""
	datasets := ""

	if len(prices) > 0 {

		sort.Slice(prices, func(i, j int) bool { return prices[i].Date.Unix() < prices[j].Date.Unix() })

		c.CurrentPrice = template.HTML(strconv.FormatFloat(prices[len(prices)-1].Price, 'f', 3, 64) + " " + prices[0].Currency)

		c.Currency = template.JS(prices[0].Currency)

		for _, price := range prices {

			if price.Type == gasType {

				switch edit {
				default:
					labels += "\"" + price.Date.Format("15:04:05") + "\","

				case "7days":
					labels += "\"" + price.Date.Format("2006-01-02 Mon") + "\","

				case "30days":
					labels += "\"" + price.Date.Format("2006-01-02") + "\","

				case "alldays":
					labels += "\"" + price.Date.Format("2006-01") + "\","
				}

				datasets += "\"" + strconv.FormatFloat(price.Price, 'f', 3, 64) + "\","
			}
		}

		if strings.HasSuffix(labels, ",") {
			labels = strings.TrimRight(labels, ",")
		}

		if strings.HasSuffix(datasets, ",") {
			datasets = strings.TrimRight(datasets, ",")
		}
	}

	c.Labels = template.JS(labels)
	c.Datasets = template.JS(datasets)

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func loadFromDB(station, gastype string, fromdate, todate time.Time) ([]price.Price, error) {

	var (
		err    error
		ret    []price.Price
		dbconn sqlite3.SQLite
		to     timeout.TimeOut
	)

	to.Start(999)

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				err = dbconn.Open(global.Config.DBFile)

			case 1:
				defer dbconn.Close()

			case 2:
				err = price.LoadConfig(&dbconn)

			case 3:
				ret, err = price.LoadList(station, gastype, fromdate, todate, dbconn.DB)
			}

			if err != nil {
				return
			}
		}
	}()

	fmt.Printf("DB: %fms\n", to.ElapsedMillis())

	return ret, err
}

func loadFromCSV(station, gastype string, fromdate, todate time.Time) ([]price.Price, error) {

	var (
		err     error
		ret     []price.Price
		plAll   pricelist.PriceList
		fileCnt uint
		wg      sync.WaitGroup
		mut     sync.Mutex
		to      timeout.TimeOut
	)

	to.Start(999)

	files, err := fileInfo.Get(global.Config.FileRoot)
	if err != nil {
		return ret, err
	}

	// shuffle the file list, to improve cpu usage
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	numCPU := runtime.NumCPU()
	chunkSize := (len(files) + numCPU - 1) / numCPU
	// split file list into chunks and start parsing
	for i := 0; i < len(files); i += chunkSize {

		end := i + chunkSize

		if end > len(files) {
			end = len(files)
		}

		wg.Add(1)

		go func(filesCore []fileInfo.FileInfo) {

			var (
				plCore  pricelist.PriceList
				cntCore uint
			)

			for _, f := range filesCore {

				if strings.HasSuffix(f.Name, ".csv") == false {
					continue
				}

				name := strings.Replace(f.Name, ".csv", "", -1)
				l := len(name)
				if l > 10 {
					name = name[:10]
				}

				fDate, err := time.Parse("2006-01-02", name)
				if err != nil {
					continue
				}

				if fDate.Before(fromdate.AddDate(0, 0, -1)) {
					continue
				}

				cntCore++

				plFile, err := csv.Import(global.Config.FileRoot+f.Name, station, gastype, fromdate)
				if err != nil {
					continue
				}

				plCore.AddList(plFile.GetList())
			}

			mut.Lock()
			plAll.AddList(plCore.GetList())
			fileCnt += cntCore
			mut.Unlock()

			wg.Done()

		}(files[i:end])
	}

	wg.Wait()

	ret = plAll.GetList()

	fmt.Printf("CSV: %fms @ %d cores (%d,%d)\n", to.ElapsedMillis(), numCPU, fileCnt, len(plAll.GetList()))

	return ret, err
}
