package global

import (
	"time"

	"github.com/dekoch/gouniversal/modules/meshFileSync/fileList"
	"github.com/dekoch/gouniversal/modules/meshFileSync/moduleConfig"
)

var (
	Config        moduleConfig.ModuleConfig
	LocalFiles    fileList.FileList
	DownloadFiles fileList.FileList
	UploadFiles   fileList.FileList

	UploadTime time.Time
)
