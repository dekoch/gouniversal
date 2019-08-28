package paratest

// only proof of concept

import (
	"fmt"

	"github.com/dekoch/gouniversal/module/paratest/einzelpara"
	"github.com/dekoch/gouniversal/module/paratest/global"
	"github.com/dekoch/gouniversal/module/paratest/paraview"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/dekoch/gouniversal/shared/timeout"
)

const dbFilePath = "data/paratest/paratest.db"

func LoadConfig() {

	global.LytAuftrag.SetTableName("auftrag")
	global.LytAuftrag.AddField("id", sqlite3.TypeINTEGER, true, true)
	global.LytAuftrag.AddField("uuid", sqlite3.TypeTEXT, false, false)
	global.LytAuftrag.AddField("name", sqlite3.TypeTEXT, false, false)
	global.LytAuftrag.AddField("typuuid", sqlite3.TypeTEXT, false, false)
	global.LytAuftrag.AddField("identuuid", sqlite3.TypeTEXT, false, false)

	global.LytParameter.SetTableName("parameter")
	global.LytParameter.AddField("id", sqlite3.TypeINTEGER, true, true)
	global.LytParameter.AddField("uuid", sqlite3.TypeTEXT, false, false)
	global.LytParameter.AddField("typ", sqlite3.TypeTEXT, false, false)
	global.LytParameter.AddField("name", sqlite3.TypeTEXT, false, false)
	global.LytParameter.AddField("prozess", sqlite3.TypeTEXT, false, false)
	global.LytParameter.AddField("parameternr", sqlite3.TypeTEXT, false, false)
	global.LytParameter.AddField("parameteruuid", sqlite3.TypeTEXT, false, false)

	global.LytEinzelPara.SetTableName("einzelpara")
	global.LytEinzelPara.AddField("id", sqlite3.TypeINTEGER, true, true)
	global.LytEinzelPara.AddField("uuid", sqlite3.TypeTEXT, false, false)
	global.LytEinzelPara.AddField("name", sqlite3.TypeTEXT, false, false)
	global.LytEinzelPara.AddField("wert", sqlite3.TypeTEXT, false, false)

	err := global.DBConn.Open(dbFilePath)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = global.DBConn.CreateTableFromLayout(global.LytAuftrag)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = global.DBConn.CreateTableFromLayout(global.LytParameter)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = global.DBConn.CreateTableFromLayout(global.LytEinzelPara)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = test(true, true, "")
	if err != nil {
		console.Log(err, "")
		return
	}
}

func Exit() {

	err := global.DBConn.Close()
	if err != nil {
		console.Log(err, "")
		return
	}
}

func test(saveall, loadall bool, loadUUID string) error {

	var err error

	if saveall {
		// save all Parameters
		var pvSave paraview.ParaView
		pvSave.NewAuftrag("Auftrag 1")
		pvSave.NewParameterSet(true, "11", "Auftrag 1 Format 1")
		pvSave.NewParameterSet(false, "11", "Auftrag 1 Format 1")
		err = pvSave.Save(global.DBConn.DB, global.DBConn.Tx)
		if err != nil {
			return err
		}
	}

	if loadall {
		// load all Parameters
		var to timeout.TimeOut
		to.Start(999)

		var pvLoad paraview.ParaView
		err = pvLoad.Load("Auftrag 1", "11", global.DBConn.DB)
		if err != nil {
			return err
		}

		fmt.Println(to.ElapsedMillis())
		fmt.Println("####")

		for _, p := range pvLoad.Typ {

			fmt.Println(p.Header.Prozess + " " + p.Header.ParameterNr + " " + p.Parameter.Name + "=" + p.Parameter.Wert)
		}

		for _, p := range pvLoad.Ident {

			fmt.Println(p.Header.Prozess + " " + p.Header.ParameterNr + " " + p.Parameter.Name + "=" + p.Parameter.Wert)
		}
	}

	if loadUUID != "" {
		// load single Parameter
		var ep einzelpara.EinzelPara
		_, err = ep.Load(loadUUID, global.DBConn.DB)
		if err != nil {
			return err
		}

		fmt.Println("UUID:" + ep.UUID + " Name:" + ep.Name + " Wert:" + ep.Wert)
	}

	return nil
}
