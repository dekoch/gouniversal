package userconfig

import (
	"os"
	"path/filepath"

	"github.com/dekoch/gouniversal/module/backup/backupitem"
)

type UserConfig struct {
	User    string
	Items   []backupitem.BackupItem
	Exclude []string
}

func (hc *UserConfig) LoadDefaults(user, fileroot string) {

	hc.User = user

	var n backupitem.BackupItem
	n.Active = true
	n.Path = "data/config/"
	hc.AddItem(n)

	n.Path = "data/lang/"
	hc.AddItem(n)

	n.Path = "data/ui/"
	hc.AddItem(n)
	// binary
	n.Path = filepath.Base(os.Args[0])
	hc.AddItem(n)

	n.Active = false
	n.Path = "data/log/"
	hc.AddItem(n)

	hc.AddExclude(fileroot)
}

func (hc *UserConfig) AddItem(item backupitem.BackupItem) {

	for i := range hc.Items {

		if hc.Items[i].Path == item.Path {

			hc.Items[i] = item
			return
		}
	}

	hc.Items = append(hc.Items, item)
}

func (hc *UserConfig) GetItems() []backupitem.BackupItem {

	return hc.Items
}

func (hc *UserConfig) AddExclude(path string) {

	for i := range hc.Exclude {

		if hc.Exclude[i] == path {
			return
		}
	}

	hc.Exclude = append(hc.Exclude, path)
}

func (hc *UserConfig) GetExclude() []string {

	return hc.Exclude
}
