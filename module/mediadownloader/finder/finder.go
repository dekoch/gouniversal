package finder

import (
	"net/url"
	"path"
	"strings"

	"github.com/dekoch/gouniversal/module/mediadownloader/finder/instagram"
	"github.com/dekoch/gouniversal/module/mediadownloader/global"
	"github.com/dekoch/gouniversal/module/mediadownloader/typemd"
)

func Find(ur, raw string) ([]typemd.DownloadFile, error) {

	var (
		err error
		ret []typemd.DownloadFile
		n   []typemd.DownloadFile
	)

	if strings.HasPrefix(ur, "https://www.instagram.com/") {

		ret, err = instagram.Find(ur, raw)
		if err != nil {
			return ret, err
		}
	}

	for _, e := range global.Config.Extension {

		n, err = findExtension(ur, raw, e)
		if err != nil {
			return ret, err
		}

		ret = append(ret, n...)
	}

	return ret, err
}

func findExtension(ur, raw, extension string) ([]typemd.DownloadFile, error) {

	var ret []typemd.DownloadFile

	raw = strings.Replace(raw, "'", "\"", -1)
	files := strings.Split(raw, "\"")

	for _, f := range files {

		paths := strings.SplitAfter(f, "\"")

		if len(paths) > 0 {

			p := paths[0]

			if strings.Contains(strings.ToLower(p), extension) {

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

				var n typemd.DownloadFile
				n.Url = strings.Replace(p, "\\u0026", "&", -1)

				name := strings.SplitAfter(path.Base(p), extension)

				if len(name) > 0 {
					n.Filename = name[0]
				}

				// check if url is already in list
				found := false

				for _, r := range ret {

					if r.Url == n.Url {
						found = true
					}
				}

				if found == false {
					ret = append(ret, n)
				}
			}
		}
	}

	return ret, nil
}
