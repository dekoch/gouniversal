package fileInfo

import (
	"time"

	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
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
func Get(path string) ([]FileInfo, error) {

	ret := []FileInfo{}

	list, err := functions.ReadDir(path, 0)
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

			fi.Size = ""
		} else {

			fi.ByteSize = l.Size()
			s := datasize.ByteSize(l.Size()).HumanReadable()
			fi.Size = s
		}

		ret = append(ret, fi)
	}

	return ret, nil
}

func GetRecursive(path string) ([]FileInfo, error) {

	ret, err := recursive(path, "")

	return ret, err
}

func recursive(parent string, path string) ([]FileInfo, error) {

	ret := []FileInfo{}

	list, err := Get(parent + path)

	for _, l := range list {

		if l.IsDir {
			rList, err := recursive(parent, path+l.Name+"/")
			if err != nil {
				empty := []FileInfo{}
				return empty, err
			}

			ret = append(ret, rList...)
		} else {
			ret = append(ret, l)
		}
	}

	return ret, err
}
