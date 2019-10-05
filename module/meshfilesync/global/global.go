package global

import (
	"github.com/dekoch/gouniversal/module/meshfilesync/filelist"
	"github.com/dekoch/gouniversal/module/meshfilesync/moduleconfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config        moduleconfig.ModuleConfig
	Lang          language.Language
	LocalFiles    filelist.FileList
	MeshFiles     filelist.FileList
	OutdatedFiles filelist.FileList
	DownloadFiles filelist.FileList
	UploadFiles   filelist.FileList
	IncomingFiles filelist.FileList

	CUploadReqStart = make(chan bool)
	CUploadStart    = make(chan bool)
)
