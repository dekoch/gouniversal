package downloader

import (
	"io/ioutil"
	"net/http"
	"strconv"

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

	return file.WriteFile(global.Config.FileRoot+f.Filename, b)
}

func DownloadTest() {

	var f typemd.DownloadFile

	for i := 1; i <= 100; i++ {

		f.Filename = strconv.Itoa(i) + ""
		f.Url = "" + strconv.Itoa(i) + ""

		Download(f)
	}
}
