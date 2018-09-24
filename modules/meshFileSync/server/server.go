package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/modules/meshFileSync/fileList"
	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/syncFile"
	"github.com/dekoch/gouniversal/modules/meshFileSync/typesMFS"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const localTest = false

var (
	fileRoot string
	fl       fileList.FileList
)

func LoadConfig() {

	fileRoot = global.Config.FileRoot

	if localTest {
		fileRoot = "input/"

		fl.SetPath(fileRoot)
		fl.SetServerID("09c8024b-bcd9-4fc5-a0e1-f2b7e63ae780")
	}
}

func Server(input typesMesh.ServerMessage) error {

	var msg typesMFS.Message

	err := json.Unmarshal(input.Message.Content, &msg)
	if err == nil {

		switch msg.Type {
		case typesMFS.MessList:

			err = list(msg, input.Sender)

		case typesMFS.MessFileUploadReq:

			err = uploadReq(msg, input.Sender)

		case typesMFS.MessFileUploadStart:

			err = uploadStart(msg, input.Sender)

		case typesMFS.MessFileUpload:

			err = upload(msg, input.Sender)
		}
	}

	return err
}

// get filelist from client
func list(input typesMFS.Message, sender serverInfo.ServerInfo) error {

	var (
		err     error
		list    []syncFile.SyncFile
		missing []syncFile.SyncFile
	)

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				// read FileList from Message
				err = json.Unmarshal(input.Content, &list)

			case 1:
				if localTest == false {

					err = global.LocalFiles.Scan()
					missing = global.LocalFiles.GetLocalMissing(list)
				} else {
					err = fl.Scan()
					missing = fl.GetLocalMissing(list)
				}

			case 2:
				global.DownloadFiles.AddList(missing)

				// update sources
				global.LocalFiles.SourceUpdateList(list)

			case 3:
				for _, lf := range global.LocalFiles.GetRemoteDeleted(list) {

					path := fileRoot + lf.Path
					if _, err := os.Stat(path); os.IsNotExist(err) == false {
						console.Output("removing ("+sender.ID+"): "+path, "meshFS")

						err = os.Remove(path)
						if err == nil {
							global.LocalFiles.MarkAsDeleted(lf.Path)
						}
					}
				}
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}

func uploadReq(input typesMFS.Message, sender serverInfo.ServerInfo) error {

	var (
		err  error
		file syncFile.SyncFile
	)

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				// read FileInfo from Message
				err = json.Unmarshal(input.Content, &file)

			case 1:
				size := datasize.ByteSize(file.Size).HumanReadable()
				console.Output("<-? ("+size+") \""+file.Path+"\" from \""+sender.ID+"\"", "meshFS")

				file.SetDestination(sender.ID)

				global.UploadFiles.Add(file)
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}

func uploadStart(input typesMFS.Message, sender serverInfo.ServerInfo) error {

	var (
		err error
		ft  typesMFS.FileTransfer
	)

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				// read FileTransfer from Message
				err = json.Unmarshal(input.Content, &ft)

			case 1:
				size := datasize.ByteSize(ft.FileInfo.Size).HumanReadable()
				console.Output("< - ("+size+") \""+ft.FileInfo.Path+"\" from \""+sender.ID+"\"", "meshFS")

				timeOut := time.Minute * time.Duration(datasize.ByteSize(ft.FileInfo.Size).MBytes())
				global.UploadTime = time.Now().Add(timeOut)
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}

// client can upload missing files to server
func upload(input typesMFS.Message, sender serverInfo.ServerInfo) error {

	var (
		err error
		ft  typesMFS.FileTransfer
		sum []byte
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				// read FileTransfer from Message
				err = json.Unmarshal(input.Content, &ft)

			case 1:
				// write file
				size := datasize.ByteSize(ft.FileInfo.Size).HumanReadable()
				console.Output("<-- ("+size+") \""+ft.FileInfo.Path+"\" from \""+sender.ID+"\"", "meshFS")
				err = file.WriteFile(fileRoot+ft.FileInfo.Path, ft.Content)

			case 2:
				// set file date
				err = os.Chtimes(fileRoot+ft.FileInfo.Path, ft.FileInfo.ModTime, ft.FileInfo.ModTime)

			case 3:
				sum, err = file.Checksum(fileRoot + ft.FileInfo.Path)

			case 4:
				if bytes.Compare(sum, ft.FileInfo.Checksum) != 0 {
					err = errors.New("error checksum")

					os.Remove(fileRoot + ft.FileInfo.Path)
				}

			case 5:
				global.DownloadFiles.Reset()
				global.UploadTime = time.Now()
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}
