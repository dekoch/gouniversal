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
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const debug = false

const dSendList time.Duration = 30 * time.Second
const dUploadReq time.Duration = dSendList
const dUpload time.Duration = 5 * time.Second

var (
	fileRoot string

	cSendListStart   = make(chan bool)
	cSendListFinish  = make(chan bool)
	cUploadReqStart  = make(chan bool)
	cUploadReqFinish = make(chan bool)
	cUploadStart     = make(chan bool)
	cUploadFinish    = make(chan bool)
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
	timerUploadReq := time.NewTimer(dUploadReq)
	timerUpload := time.NewTimer(dUpload)

	for {
		select {
		case <-timerSendList.C:
			timerSendList.Stop()
			cSendListStart <- true

		case <-cSendListFinish:
			timerSendList.Reset(dSendList)

		case <-timerUploadReq.C:
			timerUploadReq.Stop()
			cUploadReqStart <- true

		case <-cUploadReqFinish:
			timerUploadReq.Reset(dUploadReq)

		case <-timerUpload.C:
			timerUpload.Stop()
			cUploadStart <- true

		case <-cUploadFinish:
			timerUpload.Reset(dUpload)
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
		<-cUploadReqStart

		if time.Since(global.UploadTime).Seconds() > 1.0 {

			var (
				err      error
				b        []byte
				server   serverinfo.ServerInfo
				serverID string

				msg typesmfs.Message
			)

			thisID := mesh.GetServerInfo().ID

			missingFiles := global.DownloadFiles.Get()

			for _, missingFile := range missingFiles {

				if missingFile.Deleted {
					continue
				}

				if datasize.ByteSize(missingFile.Size).MBytes() > global.Config.GetMaxFileSize() {
					continue
				}

				func() {

					err = nil

					for i := 0; i <= 5; i++ {

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

						case 5:
							global.DownloadFiles.Delete(missingFile.Path)
						}

						if err != nil {
							console.Log(err, "")
							return
						}
					}
				}()
			}
		}

		cUploadReqFinish <- true
	}
}

// upload missing files to servers
func uploadFiles() {

	for {
		<-cUploadStart

		var (
			err      error
			b        []byte
			server   serverinfo.ServerInfo
			serverID string

			msg typesmfs.Message
			ft  typesmfs.FileTransfer
		)

		missingFiles := global.UploadFiles.Get()

		for _, missingFile := range missingFiles {

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
						ft.FileInfo = missingFile
						err = ft.FileInfo.Update(fileRoot)

					case 3:
						msg.Version = 1.0
						msg.Type = typesmfs.MessFileUploadStart
						msg.Content, err = json.Marshal(ft)

					case 4:
						b, err = json.Marshal(msg)

					case 5:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("  > ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						err = mesh.NewMessage(server, typemesh.MessFileSync, 1.0, b)

					case 6:
						ft.Content, err = file.ReadFile(fileRoot + ft.FileInfo.Path)

					case 7:
						msg.Type = typesmfs.MessFileUpload
						msg.Content, err = json.Marshal(ft)

					case 8:
						b, err = json.Marshal(msg)

					case 9:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("- > ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						err = mesh.NewMessage(server, typemesh.MessFileSync, 1.0, b)

						if err == nil {
							console.Output("--> ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")
						}
					}

					if err != nil {
						console.Log(err, "")
						return
					}
				}
			}()
		}

		global.UploadFiles.Reset()

		cUploadFinish <- true
	}
}
