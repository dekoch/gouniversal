package finder

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/dekoch/gouniversal/modules/mediadownloader/global"
	"github.com/dekoch/gouniversal/modules/mediadownloader/typesMD"
)

func Find(ur, raw string) ([]typesMD.DownloadFile, error) {

	var (
		ret []typesMD.DownloadFile
		n   []typesMD.DownloadFile
		err error
	)

	if strings.HasPrefix(ur, "https://www.instagram.com/") {

		ret, err = findOnInstagram(raw)
		if err == nil {
			return ret, err
		}
	}

	for _, e := range global.Config.Extension {

		n, err = findExtension(ur, raw, e)
		ret = append(ret, n...)
	}

	return ret, err
}

func findExtension(ur, raw, extension string) ([]typesMD.DownloadFile, error) {

	var ret []typesMD.DownloadFile

	files := strings.Split(raw, "href=\"")

	for _, f := range files {

		paths := strings.SplitAfter(f, "\"")

		for _, p := range paths {

			if strings.Contains(strings.ToLower(p), extension) && strings.Contains(p, "<") == false {

				p = strings.Replace(p, "\"", "", -1)

				if strings.HasPrefix(p, "//") {

					p = strings.Replace(p, "//", "", -1)

					u, err := url.Parse(ur)
					if err == nil {
						if u.Scheme == "http" {
							p = "http://" + p
						} else {
							p = "https://" + p
						}
					}

				} else if strings.HasPrefix(p, "http://") == false &&
					strings.HasPrefix(p, "https://") == false {

					if strings.HasPrefix(p, "/") == false {
						p = "/" + p
					}

					u, err := url.Parse(ur)
					if err == nil {
						if u.Scheme == "http" {
							p = "http://" + u.Host + p
						} else {
							p = "https://" + u.Host + p
						}
					}
				}

				var n typesMD.DownloadFile
				n.Url = p

				name := strings.SplitAfter(path.Base(p), extension)

				if len(name) > 0 {
					n.Filename = name[0]
				}

				ret = append(ret, n)
			}
		}
	}

	return ret, nil
}

func findOnInstagram(raw string) ([]typesMD.DownloadFile, error) {

	var ret []typesMD.DownloadFile

	found, file, err := instagramMeta(raw, "video", ".mp4")
	if err == nil {
		if found {
			ret = append(ret, file)
		}
	}

	if found == false {

		found, file, err = instagramMeta(raw, "image", ".jpg")
		if err == nil {
			if found {
				ret = append(ret, file)
			}
		}
	}

	return ret, err
}

func instagramMeta(raw, property, extension string) (bool, typesMD.DownloadFile, error) {

	var ret typesMD.DownloadFile

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

				return true, ret, nil
			}
		}
	}

	return false, ret, errors.New("not found")
}
