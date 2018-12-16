package global

import (
	"time"

	"github.com/dekoch/gouniversal/modules/meshfilesync/filelist"
	"github.com/dekoch/gouniversal/modules/meshfilesync/moduleconfig"
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

	UploadTime time.Time
)
