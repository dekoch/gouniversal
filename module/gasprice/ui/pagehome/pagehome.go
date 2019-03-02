package pagehome

import (
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/csv"
	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/lang"
	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
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
	}

	labels := ""
	datasets := ""

	fileName := time.Now().Format("2006")
	fileName += ".csv"

	pl, err := csv.Import(global.Config.FileRoot+fileName, id, gasType, from)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	prices := pl.GetList()

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
