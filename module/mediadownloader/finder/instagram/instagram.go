package instagram

import (
	"path"
	"strings"

	"github.com/dekoch/gouniversal/module/mediadownloader/typemd"
)

func Find(ur, raw string) ([]typemd.DownloadFile, error) {

	var ret []typemd.DownloadFile

	file, err := instagramMeta(raw, "video", ".mp4")
	if err != nil {
		return ret, err
	}

	ret = append(ret, file)

	files, err := instagramJSON(raw, ".jpg")
	if err != nil {
		return ret, err
	}

	ret = append(ret, files...)

	return ret, err
}

func instagramMeta(raw, property, extension string) (typemd.DownloadFile, error) {

	var ret typemd.DownloadFile

	files := strings.Split(raw, "<meta property=\"og:"+property+"\" content=\"")

	if len(files) == 2 {

		paths := strings.SplitAfter(files[1], "\"")

		if len(paths) > 1 {

			if strings.Contains(strings.ToLower(paths[0]), extension) {

				ret.Url = strings.Replace(paths[0], "\"", "", -1)

				name := strings.SplitAfter(path.Base(ret.Url), extension)

				if len(name) > 0 {
					ret.Filename = name[0]
				}
			}
		}
	}

	return ret, nil
}

func instagramJSON(raw, extension string) ([]typemd.DownloadFile, error) {

	var (
		ret []typemd.DownloadFile
		n   typemd.DownloadFile
	)

	files := strings.Split(raw, "\"display_url\":\"")

	for _, file := range files {

		if strings.HasPrefix(file, "http") == false {
			continue
		}

		u := strings.Split(file, "\",")

		if len(u) >= 1 {

			n.Url = strings.Replace(u[0], "\"", "", -1)
			n.Url = strings.Replace(n.Url, "\\u0026", "&", -1)

			name := strings.SplitAfter(path.Base(n.Url), extension)

			if len(name) > 0 {

				// check if url is already in list
				found := false

				for _, r := range ret {

					if r.Url == n.Url {
						found = true
					}
				}

				if found == false {
					n.Filename = name[0]

					ret = append(ret, n)
				}
			}
		}
	}

	return ret, nil
}
