package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/module/meshfilesync/filelist"
	"github.com/dekoch/gouniversal/module/meshfilesync/global"
	"github.com/dekoch/gouniversal/module/meshfilesync/syncfile"
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const localTest = false

var (
	fileRoot string
	tempRoot string
	fl       filelist.FileList
)

func LoadConfig() {

	fileRoot = global.Config.FileRoot
	tempRoot = global.Config.TempRoot

	if localTest {
		fileRoot = "input/"

		fl.SetPath(fileRoot)
		fl.SetServerID("09c8024b-bcd9-4fc5-a0e1-f2b7e63ae780")
	}
}

func Server(input typemesh.ServerMessage) error {

	var msg typesmfs.Message

	err := json.Unmarshal(input.Message.Content, &msg)
	if err == nil {

		switch msg.Type {
		case typesmfs.MessList:

			err = list(msg, input.Sender)

		case typesmfs.MessFileUploadReq:

			err = uploadReq(msg, input.Sender)

		case typesmfs.MessFileUploadStart:

			err = uploadStart(msg, input.Sender)

		case typesmfs.MessFileUpload:

			err = upload(msg, input.Sender)
		}
	}

	return err
}

// get filelist from client
func list(input typesmfs.Message, sender serverinfo.ServerInfo) error {

	var (
		err      error
		list     []syncfile.SyncFile
		missing  []syncfile.SyncFile
		outdated []syncfile.SyncFile
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
					outdated = global.LocalFiles.GetLocalOutdated(list)
				} else {
					err = fl.Scan()
					missing = fl.GetLocalMissing(list)
					outdated = fl.GetLocalOutdated(list)
				}

			case 2:
				if global.Config.GetAutoAdd() {

					global.DownloadFiles.AddList(missing)
				} else {

					global.MeshFiles.AddList(missing)
				}

				if global.Config.GetAutoUpdate() {

					global.DownloadFiles.AddList(outdated)
				} else {

					global.OutdatedFiles.AddList(outdated)
				}

				// update sources
				global.LocalFiles.SourceUpdateList(list)
				global.DownloadFiles.SourceUpdateList(list)

			case 3:
				if global.Config.GetAutoDelete() {

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
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}

func uploadReq(input typesmfs.Message, sender serverinfo.ServerInfo) error {

	var (
		err  error
		file syncfile.SyncFile
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

func uploadStart(input typesmfs.Message, sender serverinfo.ServerInfo) error {

	var (
		err error
		ft  typesmfs.FileTransfer
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
func upload(input typesmfs.Message, sender serverinfo.ServerInfo) error {

	var (
		err error
		ft  typesmfs.FileTransfer
		sum []byte
	)

	func() {

		for i := 0; i <= 6; i++ {

			switch i {
			case 0:
				// read FileTransfer from Message
				err = json.Unmarshal(input.Content, &ft)

			case 1:
				// write file to temp dir
				size := datasize.ByteSize(ft.FileInfo.Size).HumanReadable()
				console.Output("<-- ("+size+") \""+ft.FileInfo.Path+"\" from \""+sender.ID+"\"", "meshFS")
				err = file.WriteFile(tempRoot+ft.FileInfo.Path, ft.Content)

			case 2:
				sum, err = file.Checksum(tempRoot + ft.FileInfo.Path)

			case 3:
				if bytes.Compare(sum, ft.FileInfo.Checksum) != 0 {
					err = errors.New("error checksum")

					os.Remove(tempRoot + ft.FileInfo.Path)
				}

			case 4:
				// set file date
				err = os.Chtimes(tempRoot+ft.FileInfo.Path, ft.FileInfo.ModTime, ft.FileInfo.ModTime)

			case 5:
				// move file from temp to file dir
				global.LocalFiles.Lock()
				err = os.Rename(tempRoot+ft.FileInfo.Path, fileRoot+ft.FileInfo.Path)
				global.LocalFiles.Unlock()

			case 6:
				global.LocalFiles.Add(ft.FileInfo)
				global.DownloadFiles.Delete(ft.FileInfo.Path)
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
