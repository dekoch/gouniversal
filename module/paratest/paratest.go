package paratest

// only proof of concept

//SELECT id, uuid, name, wert FROM einzelpara WHERE uuid="5562cbd3-4000-4429-9405-d89d9baceb69" ORDER BY id DESC LIMIT 0, 1
//SELECT id, uuid, name, wert FROM einzelpara WHERE uuid="5562cbd3-4000-4429-9405-d89d9baceb69" ORDER BY id ASC

import (
	"fmt"

	"github.com/dekoch/gouniversal/module/paratest/auftrag"
	"github.com/dekoch/gouniversal/module/paratest/einzelpara"
	"github.com/dekoch/gouniversal/module/paratest/global"
	"github.com/dekoch/gouniversal/module/paratest/parameter"
	"github.com/dekoch/gouniversal/module/paratest/paraview"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/timeout"
)

const dbFilePath = "data/paratest/paratest.db"

func LoadConfig() {

	err := global.DBConn.Open(dbFilePath)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = auftrag.LoadConfig(&global.DBConn)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = parameter.LoadConfig(&global.DBConn)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = einzelpara.LoadConfig(&global.DBConn)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = test(false, true, "")
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
