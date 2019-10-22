package pagegroupedit

import (
	"errors"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/groupconfig"
	"github.com/dekoch/gouniversal/program/groupmanagement"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"

	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:Group:Edit", page.Lang.Settings.Group.GroupEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	var err error

	type content struct {
		Lang     lang.SettingsGroupEdit
		Group    groupconfig.Group
		CmbState template.HTML
		Pages    template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Group.GroupEdit

	// Form input
	id := nav.Parameter("UUID")

	if id == "new" {

		id, err = newGroup()
		if err == nil {
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
			return
		}
	}

	switch r.FormValue("edit") {
	case "apply":
		err = editGroup(r, id)
		if err == nil {
			nav.RedirectPath("Program:Settings:Group:List", false)
			return
		}

	case "delete":
		err = deleteGroup(id)
		if err == nil {
			nav.RedirectPath("Program:Settings:Group:List", false)
			return
		}
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	// copy group from array
	c.Group, err = global.GroupConfig.Get(id)
	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	// combobox State
	cmbState := "<select name=\"state\">"
	statetext := ""

	for i := 1; i <= 2; i++ {

		switch i {
		case 1:
			statetext = page.Lang.Settings.Group.GroupEdit.States.Active
		case 2:
			statetext = page.Lang.Settings.Group.GroupEdit.States.Inactive
		}

		cmbState += "<option value=\"" + strconv.Itoa(i) + "\""

		if c.Group.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	c.CmbState = template.HTML(cmbState)

	// list of pages
	pagelist := ""
	pages := nav.Sitemap.PageList()

	sort.Slice(pages, func(i, j int) bool { return pages[i] < pages[j] })

	for _, page := range pages {

		allowed := groupmanagement.IsPageAllowed(page, c.Group.UUID, false)

		if allowed {
			pagelist += "<tr class=\"table-success\">"
		} else {
			pagelist += "<tr>"
		}

		pagelist += "<td><input type=\"checkbox\" name=\"selectedpages\" value=\"" + page + "\""

		if allowed {
			pagelist += " checked"
		}

		pagelist += "></td>"

		pagelist += "<td class=\"text-left\">" + page + "</td>"
		pagelist += "</tr>"
	}

	c.Pages = template.HTML(pagelist)

	// display group
	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"settings/groupedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newGroup() (string, error) {

	u := uuid.Must(uuid.NewRandom())

	var newGroup groupconfig.Group
	newGroup.UUID = u.String()
	newGroup.Name = u.String()
	newGroup.State = 1 // active

	global.GroupConfig.Add(newGroup)

	err := global.GroupConfig.SaveConfig()

	return u.String(), err
}

func editGroup(r *http.Request, uid string) error {

	var (
		err      error
		name     string
		strState string
		intState int
		comment  string
		g        groupconfig.Group
	)

	func() {

		for i := 0; i <= 9; i++ {

			switch i {
			case 0:
				name, err = functions.CheckFormInput("name", r)

			case 1:
				strState, err = functions.CheckFormInput("state", r)

			case 2:
				comment, err = functions.CheckFormInput("comment", r)

			case 3:
				// check input
				if functions.IsEmpty(name) ||
					functions.IsEmpty(strState) {

					err = errors.New("bad input")
				}

			case 4:
				intState, err = strconv.Atoi(strState)

			case 5:
				if intState < 1 ||
					intState > 2 {

					err = errors.New("bad input")
				}

			case 6:
				g, err = global.GroupConfig.Get(uid)

			case 7:
				g.Name = name
				g.State = intState
				g.Comment = comment
				g.AllowedPages = r.Form["selectedpages"]

			case 8:
				err = global.GroupConfig.Edit(g)

			case 9:
				err = global.GroupConfig.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteGroup(uid string) error {

	global.GroupConfig.Delete(uid)

	return global.GroupConfig.SaveConfig()
}
