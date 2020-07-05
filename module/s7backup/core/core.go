package core

import (
	"errors"
	"os"
	"time"

	"github.com/dekoch/gouniversal/module/s7backup/global"
	"github.com/dekoch/gouniversal/module/s7backup/moduleconfig/scheduleconfig"
	"github.com/dekoch/gouniversal/module/s7backup/plcconfig"
	"github.com/dekoch/gouniversal/module/s7backup/s7"
	"github.com/dekoch/gouniversal/shared/console"
)

var (
	chanCheckFinished = make(chan bool)
)

func LoadConfig() error {

	go job()

	return nil
}

func job() {

	intvlCheck := 10 * time.Second
	tCheck := time.NewTimer(intvlCheck)

	for {

		select {
		case <-tCheck.C:
			tCheck.Stop()

			go check()

		case <-chanCheckFinished:
			tCheck.Reset(intvlCheck)
		}
	}
}

func check() {

	defer func() {
		chanCheckFinished <- true
	}()

	today := int(time.Now().Weekday())

	for _, bs := range global.Config.Schedule.GetList() {

		if len(bs.Day) <= today {
			continue
		}

		if bs.Day[today] == false {
			continue
		}

		if bs.Backup.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			continue
		}

		if time.Since(bs.GetChecked()) < time.Duration(15*time.Minute) {
			continue
		}

		err := backupPLC(&bs)
		if err != nil {
			console.Log(err, "S7Backup Check()")
		}
	}
}

func backupPLC(bs *scheduleconfig.BackupSchedule) error {

	var (
		err error
		pc  plcconfig.PLCConfig
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = global.Config.Schedule.SetChecked(bs.UUID, time.Now())

			case 1:
				if _, errStat := os.Stat(global.Config.GetFileRoot() + bs.PLC + ".sqlite3"); os.IsNotExist(errStat) {

					err = global.Config.Schedule.Delete(bs.UUID)
					if err != nil {
						return
					}

					err = global.Config.SaveConfig()
					return
				}

			case 2:
				_, err = pc.Load(global.Config.GetFileRoot(), bs.PLC)

			case 3:
				err = s7.BackupDB(bs.DB, &pc)

			case 4:
				err = global.Config.Schedule.SetBackup(bs.UUID, time.Now())

			case 5:
				err = global.Config.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	if err != nil {
		err = errors.New(pc.Name + " " + err.Error())
	}

	return err
}
