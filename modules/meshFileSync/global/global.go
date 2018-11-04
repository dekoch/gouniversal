package global

import (
	"time"

	"github.com/dekoch/gouniversal/modules/meshFileSync/fileList"
	"github.com/dekoch/gouniversal/modules/meshFileSync/moduleConfig"
	"github.com/dekoch/gouniversal/shared/language"
)

var (
	Config        moduleConfig.ModuleConfig
	Lang          language.Language
	LocalFiles    fileList.FileList
	MeshFiles     fileList.FileList
	OutdatedFiles fileList.FileList
	DownloadFiles fileList.FileList
	UploadFiles   fileList.FileList

	UploadTime time.Time
)
