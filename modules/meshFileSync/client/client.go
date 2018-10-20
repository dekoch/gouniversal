package client

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/typesMFS"
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

			b, err := json.Marshal(global.LocalFiles.Get())
			if err != nil {
				fmt.Println(err)
			} else {

				var msg typesMFS.Message
				msg.Type = typesMFS.MessList
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

						go func(s serverInfo.ServerInfo, jsonReq []byte) {

							defer wg.Done()

							mesh.NewMessage(s, typesMesh.MessFileSync, 1.0, jsonReq)
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
				server   serverInfo.ServerInfo
				serverID string

				msg typesMFS.Message
			)

			thisID := mesh.GetServerInfo().ID

			missingFiles := global.DownloadFiles.Get()

			for _, missingFile := range missingFiles {

				if missingFile.Deleted {
					continue
				}

				if datasize.ByteSize(missingFile.Size).MBytes() > global.Config.MaxFileSize {
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
							msg.Type = typesMFS.MessFileUploadReq
							msg.Version = 1.0
							msg.Content, err = json.Marshal(missingFile)

						case 3:
							b, err = json.Marshal(msg)

						case 4:
							size := datasize.ByteSize(missingFile.Size).HumanReadable()
							console.Output("?-> ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

							err = mesh.NewMessage(server, typesMesh.MessFileSync, 1.0, b)

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
			server   serverInfo.ServerInfo
			serverID string

			msg typesMFS.Message
			ft  typesMFS.FileTransfer
		)

		missingFiles := global.UploadFiles.Get()

		for _, missingFile := range missingFiles {

			if _, err := os.Stat(fileRoot + missingFile.Path); os.IsNotExist(err) {
				continue
			}

			if datasize.ByteSize(missingFile.Size).MBytes() > global.Config.MaxFileSize {
				continue
			}

			func() {

				err = nil

				for i := 0; i <= 8; i++ {

					switch i {
					case 0:
						serverID = missingFile.GetDestination()

					case 1:
						server, err = mesh.GetServerWithID(serverID)

					case 2:
						ft.FileInfo = missingFile

						msg.Version = 1.0
						msg.Type = typesMFS.MessFileUploadStart
						msg.Content, err = json.Marshal(ft)

					case 3:
						b, err = json.Marshal(msg)

					case 4:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("  > ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						err = mesh.NewMessage(server, typesMesh.MessFileSync, 1.0, b)

					case 5:
						ft.Content, err = file.ReadFile(fileRoot + ft.FileInfo.Path)

					case 6:
						msg.Type = typesMFS.MessFileUpload
						msg.Content, err = json.Marshal(ft)

					case 7:
						b, err = json.Marshal(msg)

					case 8:
						size := datasize.ByteSize(missingFile.Size).HumanReadable()
						console.Output("- > ("+size+") \""+missingFile.Path+"\" to \""+serverID+"\"", "meshFS")

						err = mesh.NewMessage(server, typesMesh.MessFileSync, 1.0, b)

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