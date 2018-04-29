package TempHumV1x0request

// http://127.0.0.1:8080/request/?id=test&key=1234

import (
	"encoding/json"
	"gouniversal/modules/openespm/app/TempHumV1x0"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/functions"
	"gouniversal/shared/io/csv"
	"strconv"
	"time"
)

type appResp struct {
	Dev typesOESPM.DefaultDevResp
}

func Request(resp *typesOESPM.Response, req *typesOESPM.Request) {

	// read data and write .csv
	ctemp := req.Values.Get("ctemp")
	ftemp := req.Values.Get("ftemp")
	humidity := req.Values.Get("humidity")

	if functions.IsEmpty(ctemp) == false &&
		functions.IsEmpty(ftemp) == false &&
		functions.IsEmpty(humidity) == false {

		t := time.Now()

		row := []string{}
		row = append(row, t.Format("2006-01-02"))
		row = append(row, t.Format("15:04:05"))
		row = append(row, strconv.Itoa(int(t.Unix())))
		row = append(row, ctemp)
		row = append(row, ftemp)
		row = append(row, humidity)

		fileDir := req.DeviceDataFolder +
			strconv.Itoa(t.Year()) + "/" +
			strconv.Itoa(int(t.Month())) + "/"

		filePath := fileDir + strconv.Itoa(t.Day()) + ".csv"

		err := csv.AddRow(filePath, row)
		if err != nil {
			resp.Err = err
		}
	}

	resp.Type = typesOESPM.JSON

	// init new device
	if functions.IsEmpty(req.Device.Config) {
		req.Device.Config = TempHumV1x0.InitDeviceConfig()
	}

	// read device config
	var c TempHumV1x0.DeviceConfig
	err := req.Device.Unmarshal(&c)
	if err != nil {
		resp.Err = err
		return
	}

	// build json response
	var r appResp
	r.Dev.Ver = 1.0
	r.Dev.Intvl = c.Dev.Interval()
	r.Dev.Ds = true

	b, err := json.Marshal(r)
	if err != nil {
		resp.Err = err
	} else {
		resp.Content = string(b[:])
	}

	// write device config
	resp.Err = req.Device.Marshal(c)
}
