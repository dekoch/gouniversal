package pagehome

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/module/instabackup/global"
	"github.com/dekoch/gouniversal/module/instabackup/instafile"
	"github.com/dekoch/gouniversal/module/instabackup/lang"
	"github.com/dekoch/gouniversal/module/instabackup/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/sqlite3"
	"github.com/dekoch/gouniversal/shared/navigation"
)

const TableName = "instafile"

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:InstaBackup:Home", page.Lang.Home.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang        lang.Home
		Gallery     template.HTML
		PreviousSet template.HTML
		NextSet     template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	var (
		err       error
		inputSet  string
		selectSet int
		files     []instafile.InstaFile
		token     string
	)

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				inputSet, err = functions.CheckFormInput("next", r)

			case 1:
				if functions.IsEmpty(inputSet) {
					selectSet = 0
				} else {
					selectSet, err = strconv.Atoi(inputSet)
				}

			case 2:
				if selectSet < 0 {
					selectSet = 0
				}

				instaID := global.Config.GetIDFromUser(nav.User.UUID)
				if len(instaID) == 0 {
					return
				}

				files, err = getFiles(instaID, selectSet, 12)

			case 3:
				token = global.Tokens.New(nav.User.UUID)
				if token == "" {
					err = errors.New("token not set")
				}
			}

			if err != nil {
				return
			}
		}
	}()

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	tbody := ""

	for _, file := range files {

		tbody += "<tr><td>"
		tbody += "<a href=\"https://www.instagram.com/" + file.UserName + "\" target=\"_blank\" rel=\"noreferrer\">" + file.UserName + "</a><br>"

		if strings.HasSuffix(file.FileID, ".jpg") {
			// image
			//tbody += "<img src=\"" + file.URL + "\" height=\"400\">"
			tbody += "<img src=\"/instabackup/req/?uuid=" + nav.User.UUID + "&token=" + token + "&file=" + file.UserID + "/" + file.UserID + "_" + file.UserName + "_" + file.FileID + "\" height=\"400\">"
		} else {
			// video
			tbody += "<video height=\"400\" controls>"
			//tbody += "<source src=\"" + file.URL + "\" type=\"video/mp4\">"
			tbody += "<source src=\"/instabackup/req/?uuid=" + nav.User.UUID + "&token=" + token + "&file=" + file.UserID + "/" + file.UserID + "_" + file.UserName + "_" + file.FileID + "\" type=\"video/mp4\">"
			tbody += "Your browser does not support the video tag."
			tbody += "</video>"
		}

		tbody += "</td></tr>"
	}

	c.Gallery = template.HTML(tbody)
	c.PreviousSet = template.HTML(strconv.Itoa(selectSet - 12))
	c.NextSet = template.HTML(strconv.Itoa(selectSet + 12))

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func getFiles(users []string, start, lastn int) ([]instafile.InstaFile, error) {

	if len(users) == 0 {
		return []instafile.InstaFile{}, errors.New("empty user list")
	}

	var (
		err    error
		ret    []instafile.InstaFile
		dbconn sqlite3.SQLite
		rows   *sql.Rows
		ids    []string
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				err = dbconn.Open(global.Config.DBFile)

			case 1:
				defer dbconn.Close()

			case 2:
				var where string

				for iu := range users {

					where += "userid='" + users[iu] + "'"

					if iu != len(users)-1 {
						where += " OR "
					}
				}

				rows, err = dbconn.DB.Query("SELECT id FROM `"+TableName+"` WHERE "+where+" ORDER BY id DESC LIMIT ?, ?", start, lastn)

			case 3:
				defer rows.Close()

			case 4:
				var id string

				for rows.Next() {

					err = rows.Scan(&id)
					if err != nil {
						return
					}

					ids = append(ids, id)
				}

			case 5:
				var (
					n     instafile.InstaFile
					found bool
				)

				for _, id := range ids {

					found, err = n.Load(id, &dbconn)
					if err != nil {
						return
					}

					if found {
						ret = append(ret, n)
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}
