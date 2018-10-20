package downloader

import (
	"io/ioutil"
	"net/http"

	"github.com/dekoch/gouniversal/modules/mediadownloader/global"
	"github.com/dekoch/gouniversal/modules/mediadownloader/typesMD"
	"github.com/dekoch/gouniversal/shared/io/file"
)

func Download(f typesMD.DownloadFile) error {

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
