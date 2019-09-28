package fileinfo

import (
	"io/ioutil"
	"time"

	"github.com/dekoch/gouniversal/shared/datasize"
)

type FileInfo struct {
	Name     string
	ByteSize int64
	Size     string
	ModTime  time.Time
	IsDir    bool
	Path     string
}

// Get returns folders and files from given path
func Get(path string, maxdepth int, withdir bool) ([]FileInfo, error) {

	return recursive(path, maxdepth, withdir, 0)
}

// recursive is a helper for Get()
func recursive(path string, maxdepth int, withdir bool, currdepth int) ([]FileInfo, error) {

	var ret []FileInfo

	list, err := ioutil.ReadDir(path)
	if err != nil {
		return ret, err
	}

	var fi FileInfo

	for _, l := range list {

		fi.Name = l.Name()
		fi.ModTime = l.ModTime()
		fi.IsDir = l.IsDir()
		fi.Path = path

		if fi.IsDir {

			fi.ByteSize = 0
			fi.Size = ""

		} else {

			fi.ByteSize = l.Size()
			s := datasize.ByteSize(l.Size()).HumanReadable()
			fi.Size = s
		}

		ret = append(ret, fi)
	}

	for _, r := range ret {

		if r.IsDir {

			if maxdepth > currdepth ||
				maxdepth < 0 {

				rList, err := recursive(path+r.Name+"/", maxdepth, withdir, currdepth+1)
				if err != nil {
					return ret, err
				}

				ret = append(ret, rList...)
			}
		}
	}

	if withdir {
		return ret, nil
	}

	var fList []FileInfo

	for _, r := range ret {

		if r.IsDir == false {
			fList = append(fList, r)
		}
	}

	return fList, nil
}
