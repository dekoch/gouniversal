package fileInfo

import (
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
)

type FileInfo struct {
	Name string
	Size string
}

func Get(path string) ([]FileInfo, []FileInfo) {

	list, _ := functions.ReadDir(path, 0)

	folders := []FileInfo{}
	files := []FileInfo{}

	var fi FileInfo

	for _, l := range list {

		fi.Name = l.Name()

		if l.IsDir() {

			fi.Size = ""
			folders = append(folders, fi)
		} else {

			s := datasize.ByteSize(l.Size()).HumanReadable()
			fi.Size = s
			files = append(files, fi)
		}
	}

	return folders, files
}
