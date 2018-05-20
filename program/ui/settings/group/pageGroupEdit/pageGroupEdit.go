package pageGroupEdit

import (
	"errors"
	"gouniversal/program/global"
	"gouniversal/program/groupConfig"
	"gouniversal/program/groupManagement"
	"gouniversal/program/lang"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "Program:Settings:Group:Edit", page.Lang.Settings.Group.GroupEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang     lang.SettingsGroupEdit
		Group    groupConfig.Group
		CmbState template.HTML
		Pages    template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Group.GroupEdit

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newGroup()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		err := editGroup(r, id)
		if err == nil {
			nav.RedirectPath("Program:Settings:Group:List", false)
		}

	} else if button == "delete" {

		err := deleteGroup(id)
		if err == nil {
			nav.RedirectPath("Program:Settings:Group:List", false)
		}
	}

	// copy group from array
	var err error
	c.Group, err = global.GroupConfig.Get(id)

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
	sm := nav.Sitemap.PageList()

	for i := 0; i < len(sm); i++ {

		pagelist += "<tr>"
		pagelist += "<td>" + sm[i] + "</td>"
		pagelist += "<td><input type=\"checkbox\" name=\"selectedpages\" value=\"" + sm[i] + "\""

		if groupManagement.IsPageAllowed(sm[i], c.Group.UUID, false) {

			pagelist += " checked"
		}
		pagelist += "></td></tr>"
	}

	c.Pages = template.HTML(pagelist)

	// display group
	p, err := functions.PageToString(global.UiConfig.File.ProgramFileRoot+"settings/groupedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newGroup() string {

	u := uuid.Must(uuid.NewRandom())

	var newGroup groupConfig.Group
	newGroup.UUID = u.String()
	newGroup.Name = u.String()
	newGroup.State = 1 // active

	global.GroupConfig.Add(newGroup)

	global.GroupConfig.SaveConfig()

	return u.String()
}

func editGroup(r *http.Request, uid string) error {

	name, _ := functions.CheckFormInput("name", r)
	state, _ := functions.CheckFormInput("state", r)
	comment, errComment := functions.CheckFormInput("comment", r)

	// check input
	if functions.IsEmpty(name) ||
		functions.IsEmpty(state) ||
		govalidator.IsNumeric(state) == false ||
		// content not required
		errComment != nil {

		return errors.New("bad input")
	}

	intState, err := strconv.Atoi(state)
	if err != nil {
		return err
	}

	selpages := r.Form["selectedpages"]

	g, err := global.GroupConfig.Get(uid)
	if err != nil {
		return err
	}

	g.Name = name
	g.State = intState
	g.Comment = comment
	g.AllowedPages = selpages

	err = global.GroupConfig.Edit(g)
	if err != nil {
		return err
	}

	return global.GroupConfig.SaveConfig()
}

func deleteGroup(uid string) error {

	global.GroupConfig.Delete(uid)

	return global.GroupConfig.SaveConfig()
}
