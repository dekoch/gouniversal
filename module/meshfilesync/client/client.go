package client

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh"
	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/module/meshfilesync/global"
	"github.com/dekoch/gouniversal/module/meshfilesync/outgoingfile"
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const debug = false

const dSendList time.Duration = 30 * time.Second

var (
	fileRoot string

	cSendListStart  = make(chan bool)
	cSendListFinish = make(chan bool)
)

func LoadConfig() {

	fileRoot = global.Config.FileRoot

	global.UploadFiles.SetPath(fileRoot)

	go sendFileList()
	go sendUploadReq()
	go uploadFiles()
	go job()
}

func job() {

	timerSendList := time.NewTimer(dSendList)

	for {
		select {
		case <-timerSendList.C:
			timerSendList.Stop()
			cSendListStart <- true

		case <-cSendListFinish:
			timerSendList.Reset(dSendList)
		}
	}
}

// get File Lists from servers
func sendFileList() {

	for {
		<-cSendListStart

		if debug {
			fmt.Println("send file list")
		}

		func() {
			err := global.LocalFiles.Scan()
			if err != nil {
				return
			}

			global.LocalFiles.SourceClean(mesh.GetServerList())

			b, err := json.Marshal(global.LocalFiles.Get())
			if err != nil {
				fmt.Println(err)
			} else {

				var msg typesmfs.Message
				msg.Type = typesmfs.MessList
				msg.Version = 1.0
				msg.Content = b

				b, err = json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
				} else {

					var wg sync.WaitGroup

					// send request to all servers
					for _, server := range mesh.GetServerList() {

						// only to other systems
						if mesh.IsLoop(server) {
							continue
						}

						if debug {
							fmt.Println(server.ID)
						}

						wg.Add(1)

						go func(s serverinfo.ServerInfo, jsonReq []byte) {

							defer wg.Done()

							mesh.NewMessage(s, typemesh.MessFileSync, 1.0, jsonReq)
						}(server, b)
					}

					wg.Wait()
				}
			}
		}()

		cSendListFinish <- true
	}
}

// send Upload Request for missing files
func sendUploadReq() {

	for {
		<-global.CUploadReqStart

		var (
			err      error
			b        []byte
			server   serverinfo.ServerInfo
			serverID string

			msg typesmfs.Message
		)

		global.IncomingFiles.ClearIncomingFiles(30.0)

		thisID := mesh.GetServerInfo().ID

		for _, missingFile := range global.DownloadFiles.Get() {

			if missingFile.Deleted {
				continue
			}

			if global.IncomingFiles.Exists(missingFile.Path) {
				continue
			}

			if datasize.ByteSize(missingFile.Size).MBytes() > global.Config.GetMaxFileSize() {
				continue
			}

			func() {

				err = nil

				for i := 0; i <= 4; i++ {

					switch i {
					case 0:
						serverID, err = missingFile.SelectSource(thisID)

					case 1:
						server, err = mesh.GetServerWithID(serverID)

					case 2:
						msg.Type = typesmfs.MessFileUploadReq
						msg.Version = 1.0
						msg.Content, err = json.Marshal(missingFile)

					case 3:
						b, err = json.Marshal(msg)

					case 4:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("?-> ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						err = mesh.NewMessage(server, typemesh.MessFileSync, 1.0, b)
					}

					if err != nil {
						console.Log(err, "")
						return
					}
				}
			}()
		}
	}
}

// upload missing files to servers
func uploadFiles() {

	for {
		<-global.CUploadStart

		var (
			err      error
			b        []byte
			server   serverinfo.ServerInfo
			serverID string
			done     bool

			msg     typesmfs.Message
			ft      typesmfs.FileTransfer
			outFile outgoingfile.OutgoingFile
		)

		for _, missingFile := range global.UploadFiles.Get() {

			if done {
				continue
			}

			if _, err := os.Stat(fileRoot + missingFile.Path); os.IsNotExist(err) {
				continue
			}

			if datasize.ByteSize(missingFile.Size).MBytes() > global.Config.GetMaxFileSize() {
				continue
			}

			func() {

				err = nil

				for i := 0; i <= 9; i++ {

					switch i {
					case 0:
						serverID = missingFile.GetDestination()

					case 1:
						server, err = mesh.GetServerWithID(serverID)

					case 2:
						_, err = missingFile.Update(fileRoot)

					case 3:
						ft.FileInfo = missingFile
						ft.Content, err = file.ReadFile(fileRoot + ft.FileInfo.Path)

					case 4:
						msg.Version = 1.0
						msg.Type = typesmfs.MessFileUpload
						msg.Content, err = json.Marshal(ft)

					case 5:
						b, err = json.Marshal(msg)

					case 6:
						outFile.Set(missingFile, server)
						err = outFile.SendNewMessage(true)

					case 7:
						err = outFile.SendContMessage(false)

					case 8:
						err = mesh.NewMessage(server, typemesh.MessFileSync, 1.0, b)

						outFile.Delete()

					case 9:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("--> ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						done = true
					}

					if err != nil {
						console.Log(err, "")
						return
					}
				}
			}()
		}

		outFile.Delete()
		global.UploadFiles.Reset()
	}
}
