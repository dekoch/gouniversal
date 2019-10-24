package pagehome

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module/backup/backupitem"
	"github.com/dekoch/gouniversal/module/backup/global"
	"github.com/dekoch/gouniversal/module/backup/lang"
	"github.com/dekoch/gouniversal/module/backup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/archive"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/io/fileinfo"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:Backup:Home", page.Lang.Home.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang  lang.Home
		Items template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	var err error

	switch r.FormValue("edit") {
	case "download":

		func() {

			for i := 0; i <= 2; i++ {

				switch i {
				case 0:
					err = editItems(nav.User.UUID, page, r)

				case 1:
					err = global.Config.SaveConfig()

				case 2:
					err = createAndServe(nav.User.UUID, page, nav, r)
				}

				if err != nil {
					return
				}
			}
		}()
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	itemsTable := ""
	items := global.Config.GetItems(nav.User.UUID)
	exclude := global.Config.GetExclude(nav.User.UUID)

	sort.Slice(items, func(i, j int) bool { return items[i].Path < items[j].Path })

	for _, item := range items {
		// checkbox
		itemsTable += "<tr>"
		itemsTable += "<td><input type=\"checkbox\" name=\"selecteditem\" value=\"" + item.Path + "\""

		if item.Active {
			itemsTable += " checked"
		}

		itemsTable += "></td>"
		// path
		itemsTable += "<td class=\"text-left\">" + item.Path + "</td>"
		// size
		dir := filepath.Dir(item.Path)
		if strings.HasSuffix(dir, "/") == false {
			dir += "/"
		}

		fi, _ := fileinfo.Get(dir, -1, true)

		var size int64

		for ifi := range fi {

			cont := false

			for iex := range exclude {

				if dir == "./" {

					if strings.HasPrefix(fi[ifi].Name, exclude[iex]) {
						cont = true
					}
				} else {

					if strings.HasPrefix(fi[ifi].Path+fi[ifi].Name, exclude[iex]) {
						cont = true
					}
				}
			}

			if cont {
				continue
			}

			if dir == "./" {

				if strings.HasPrefix(fi[ifi].Name, item.Path) {
					size += fi[ifi].ByteSize
				}
			} else {
				size += fi[ifi].ByteSize
			}
		}

		itemsTable += "<td class=\"text-left\">" + datasize.ByteSize(size).HumanReadable() + "</td>"
		itemsTable += "</tr>"
	}

	c.Items = template.HTML(itemsTable)

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func editItems(user string, page *typemd.Page, r *http.Request) error {

	var (
		err           error
		selectedItems []string
	)

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				selectedItems = r.Form["selecteditem"]

				if len(selectedItems) == 0 {
					err = errors.New(page.Lang.Home.NoItemSelected)
				}

			case 1:
				for _, item := range global.Config.GetItems(user) {

					item.Active = false

					for _, sitem := range selectedItems {

						if item.Path == sitem {
							item.Active = true
						}
					}

					global.Config.AddItem(user, item)
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func createAndServe(user string, page *typemd.Page, nav *navigation.Navigation, r *http.Request) error {

	var (
		err   error
		b     []byte
		items []backupitem.BackupItem
	)

	now := time.Now().Format("20060102_150405")
	name := now + "_backup.zip"
	tempDir := global.Config.FileRoot + now + "_backup/"
	targetFile := global.Config.FileRoot + name

	for i := 0; i <= 6; i++ {

		switch i {
		case 0:
			for _, item := range global.Config.GetItems(user) {

				if item.Active == false {
					continue
				}

				if _, err = os.Stat(item.Path); os.IsNotExist(err) {
					return err
				}

				items = append(items, item)
			}

			if len(items) == 0 {
				err = errors.New(page.Lang.Home.NoItemSelected)
			}

		case 1:
			// archive items
			for _, item := range items {

				err = archive.Zipit(item.Path, tempDir+filepath.Base(item.Path)+".zip", global.Config.GetExclude(user))
				if err != nil {
					return err
				}
			}

		case 2:
			err = archive.Zipit(tempDir, targetFile, []string{})

		case 3:
			// read archive
			b, err = file.ReadFile(targetFile)

		case 4:
			// tell the browser the returned content should be downloaded
			nav.ResponseWriter.Header().Add("Content-Disposition", "attachment; filename="+name)

			var rs io.ReadSeeker
			rs = bytes.NewReader(b)
			http.ServeContent(nav.ResponseWriter, r, name, time.Now(), rs)

		case 5:
			err = file.Remove(tempDir)

		case 6:
			err = file.Remove(targetFile)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
