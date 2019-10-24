package core

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dekoch/gouniversal/module/instabackup/global"
	"github.com/dekoch/gouniversal/module/instabackup/instafile"
	"github.com/dekoch/gouniversal/module/instabackup/instaquery"
	"github.com/dekoch/gouniversal/module/instabackup/instaresp"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

type downloadFile struct {
	FileName string
	Url      string
}

var (
	chanBackupFinished = make(chan bool)
)

func LoadConfig() {

	var dbconn sqlite3.SQLite

	err := dbconn.Open(global.Config.DBFile)
	if err != nil {
		console.Log(err, "instabackup")
		return
	}

	defer dbconn.Close()

	err = instafile.LoadConfig(&dbconn)
	if err != nil {
		console.Log(err, "instabackup")
		return
	}

	err = instaresp.LoadConfig(&dbconn)
	if err != nil {
		console.Log(err, "instabackup")
		return
	}

	go job()

	if global.Config.UpdInterv == -1 {
		go backup()
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
				timer.Stop()
				go backup()

			case <-chanBackupFinished:
				intvl = global.Config.GetUpdInterval()
				timer.Reset(intvl)
			}
		} else {
			// wait until enabled
			time.Sleep(1 * time.Minute)
			intvl = global.Config.GetUpdInterval()

			if intvl > 0 {
				timer.Reset(intvl)
			}
		}
	}
}

func backup() {

	defer func() {
		chanBackupFinished <- true
	}()

	err := global.Config.LoadConfig()
	if err != nil {
		console.Log(err, "instabackup")
		return
	}

	for _, user := range global.Config.GetUserList() {

		err = backupUser(user)
		if err != nil {
			console.Log(err, "instabackup")
			continue
		}
	}

	err = global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "instabackup")
		return
	}
}

func backupUser(userid string) error {

	var (
		err       error
		hash      string
		b         []byte
		ir        instaresp.InstaResp
		downloads []downloadFile
		files     []instafile.InstaFile
		dbconn    sqlite3.SQLite
		userName  string
	)

	func() {

		for i := 0; i <= 11; i++ {

			switch i {
			case 0:
				// check input
				if functions.IsEmpty(userid) {
					err = errors.New("invalid userid")
				}

			case 1:
				hash, err = global.Config.Hashes.GetHash()

			case 2:
				err = dbconn.Open(global.Config.DBFile)

			case 3:
				defer dbconn.Close()

			case 4:
				ir.UserID = userid
				_, err = ir.Load(&dbconn)

			case 5:
				var (
					iq          instaquery.InstaQuery
					hashExpired bool
				)

				iq.QueryHash = hash
				iq.Variables.ID = ir.UserID
				iq.Variables.First = 50

				downloads, userName, hashExpired, err = getFiles(iq, &ir)

				if hashExpired {
					global.Config.Hashes.SetAsExpired(hash, global.Config.GetHashReset())
				}

			case 6:
				var n instafile.InstaFile
				t := time.Now()
				path := global.Config.FileRoot

				for _, f := range downloads {

					n.UserID = userid
					n.UserName = userName
					n.Added = t
					n.FileID = f.FileName
					n.URL = f.Url

					found, err := n.Exists(&dbconn)
					if err != nil || found {
						continue
					}

					b, err = download(f)
					if err != nil {
						continue
					}

					filePath := userid + "/" + userid + "_" + userName + "_" + f.FileName

					console.Output("writing file: "+filePath, "instabackup")

					err = file.WriteFile(path+filePath, b)
					if err != nil {
						continue
					}

					files = append(files, n)
				}

			case 7:
				dbconn.Tx, err = dbconn.DB.Begin()

			case 8:
				defer func() {
					if err != nil {
						dbconn.Tx.Rollback()
					}
				}()

			case 9:
				ir.UserName = userName
				ir.Checked = time.Now()
				ir.Save(dbconn.Tx)

			case 10:
				// save
				for _, f := range files {

					err = f.Save(dbconn.Tx)
					if err != nil {
						return
					}
				}

			case 11:
				err = dbconn.Tx.Commit()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func getFiles(iq instaquery.InstaQuery, ir *instaresp.InstaResp) ([]downloadFile, string, bool, error) {

	var (
		err         error
		files       []downloadFile
		userName    string
		hashExpired bool
		b           []byte
	)

	func() {

		for {

			for i := 0; i <= 5; i++ {

				switch i {
				case 0:
					b, err = iq.SendQuery()

				case 1:
					err = ir.Response.Unmarshal(b)

				case 2:
					if ir.Response.Status != "ok" {
						hashExpired = true
						return
					}

				case 3:
					if len(ir.Response.Data.User.Eottm.Edges) >= 1 {

						if iq.Variables.ID == ir.Response.Data.User.Eottm.Edges[0].Node.Owner.ID {
							userName = ir.Response.Data.User.Eottm.Edges[0].Node.Owner.UserName
						}
					}

				case 4:
					var n downloadFile

					for _, f := range ir.GetFiles() {

						if f.IsVideo {

							n.FileName = f.FileID + ".mp4"
							n.Url = f.VideoURL
						} else {

							n.FileName = f.FileID + ".jpg"
							n.Url = f.DisplayURL
						}
						// check if filename is already in array
						found := false

						for ifi := range files {

							if files[ifi].FileName == n.FileName {
								found = true
							}
						}

						if found {
							continue
						}

						files = append(files, n)
					}

				case 5:
					iq.Variables.After = ir.Response.Data.User.Eottm.PageInfo.EndCursor

					if ir.Response.Data.User.Eottm.PageInfo.HasNextPage == false {
						return
					}
				}

				if err != nil {
					return
				}
			}
		}
	}()

	return files, userName, hashExpired, err
}

func download(f downloadFile) ([]byte, error) {

	resp, err := http.Get(f.Url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
