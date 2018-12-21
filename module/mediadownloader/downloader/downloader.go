package downloader

import (
	"io/ioutil"
	"net/http"

	"github.com/dekoch/gouniversal/module/mediadownloader/global"
	"github.com/dekoch/gouniversal/module/mediadownloader/typemd"
	"github.com/dekoch/gouniversal/shared/io/file"
)

func Download(f typemd.DownloadFile) error {

	resp, err := http.Get(f.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	file.WriteFile(global.Config.FileRoot+f.Filename, b)

	return nil
}
