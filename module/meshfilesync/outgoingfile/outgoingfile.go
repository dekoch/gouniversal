package outgoingfile

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh"
	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/module/meshfilesync/syncfile"
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/datasize"
)

const dSend time.Duration = 10 * time.Second

type OutgoingFile struct {
	mut       sync.Mutex
	isSet     bool
	syncFile  syncfile.SyncFile
	server    serverinfo.ServerInfo
	cFinish   chan bool
	cFinished chan bool
}

func (of *OutgoingFile) Set(sf syncfile.SyncFile, server serverinfo.ServerInfo) {

	of.mut.Lock()
	defer of.mut.Unlock()

	of.isSet = true
	of.syncFile = sf
	of.server = server

	of.cFinish = make(chan bool)
	of.cFinished = make(chan bool)
}

func (of *OutgoingFile) Delete() {

	if of.IsSet() == false {
		return
	}

	of.cFinish <- true
	<-of.cFinished

	of.mut.Lock()
	defer of.mut.Unlock()

	of.isSet = false
}

func (of *OutgoingFile) IsSet() bool {

	of.mut.Lock()
	defer of.mut.Unlock()

	return of.isSet
}

func (of *OutgoingFile) SendNewMessage(start bool) error {

	of.mut.Lock()
	defer of.mut.Unlock()

	if of.isSet == false {
		return errors.New("parameter not set")
	}

	var (
		err error
		b   []byte
		msg typesmfs.Message
		ft  typesmfs.FileTransfer
	)

	for i := 0; i <= 3; i++ {

		switch i {
		case 0:
			ft.FileInfo = of.syncFile

		case 1:
			msg.Version = 1.0
			msg.Type = typesmfs.MessFileUploadStart
			msg.Content, err = json.Marshal(ft)

		case 2:
			b, err = json.Marshal(msg)

		case 3:
			size := datasize.ByteSize(of.syncFile.Size).HumanReadable()

			if start {
				console.Output("  > ("+size+") \""+of.syncFile.Path+"\" to \""+of.server.ID+"\"", "meshFS")
			} else {
				console.Output("- > ("+size+") \""+of.syncFile.Path+"\" to \""+of.server.ID+"\"", "meshFS")
			}

			err = mesh.NewMessage(of.server, typemesh.MessFileSync, 1.0, b)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (of *OutgoingFile) SendContMessage(start bool) error {

	of.mut.Lock()
	defer of.mut.Unlock()

	if of.isSet == false {
		return errors.New("parameter not set")
	}

	go func(s bool) {

		timerSend := time.NewTimer(dSend)

		for {

			select {
			case <-of.cFinish:
				of.cFinished <- true
				return

			case <-timerSend.C:
				timerSend.Stop()

				of.SendNewMessage(s)

				timerSend.Reset(dSend)
			}
		}
	}(start)

	return nil
}
