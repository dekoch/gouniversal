package pageGroupEdit

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/groupManagement"
	"gouniversal/program/lang"
	"gouniversal/program/programTypes"
	"gouniversal/program/ui/uifunc"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:Group:Edit", page.Lang.Settings.Group.GroupEdit.Title)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type groupedit struct {
		Lang     lang.SettingsGroupEdit
		Group    programTypes.Group
		CmbState template.HTML
		Pages    template.HTML
	}
	var ge groupedit

	ge.Lang = page.Lang.Settings.Group.GroupEdit

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newGroup()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		editGroup(r, id)

		nav.RedirectPath("Program:Settings:Group:List", false)

	} else if button == "delete" {

		deleteGroup(id)

		nav.RedirectPath("Program:Settings:Group:List", false)
	}

	// copy group from array
	global.GroupConfig.Mut.Lock()
	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		if id == global.GroupConfig.File.Group[i].UUID {

			ge.Group = global.GroupConfig.File.Group[i]
		}
	}
	global.GroupConfig.Mut.Unlock()

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

		if ge.Group.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	ge.CmbState = template.HTML(cmbState)

	// list of pages
	pagelist := ""
	sm := nav.Sitemap.PageList()

	for i := 0; i < len(sm); i++ {

		pagelist += "<tr>"
		pagelist += "<td>" + sm[i] + "</td>"
		pagelist += "<td><input type=\"checkbox\" name=\"selectedpages\" value=\"" + sm[i] + "\""

		if groupManagement.IsPageAllowed(sm[i], ge.Group.UUID, false) {

			pagelist += " checked"
		}
		pagelist += "></td></tr>"
	}

	ge.Pages = template.HTML(pagelist)

	// display group
	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/settings/groupedit.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, ge)
}

func newGroup() string {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())

	newgroup := make([]programTypes.Group, 1)
	newgroup[0].UUID = u.String()
	newgroup[0].Name = u.String()
	newgroup[0].State = 1 // active

	global.GroupConfig.File.Group = append(newgroup, global.GroupConfig.File.Group...)

	groupManagement.SaveGroup(global.GroupConfig.File)

	return u.String()
}

func editGroup(r *http.Request, u string) {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	name := uifunc.CheckFormInput("name", r)
	state := uifunc.CheckFormInput("state", r)

	if uifunc.CheckInput(name, uifunc.STRING) &&
		uifunc.CheckInput(state, uifunc.INT) {

		for i := 0; i < len(global.GroupConfig.File.Group); i++ {

			if u == global.GroupConfig.File.Group[i].UUID {

				intState, err := strconv.Atoi(state)

				if err == nil {
					comment := uifunc.CheckFormInput("comment", r)

					selpages := r.Form["selectedpages"]

					global.GroupConfig.File.Group[i].Name = name
					global.GroupConfig.File.Group[i].State = intState
					global.GroupConfig.File.Group[i].Comment = comment

					global.GroupConfig.File.Group[i].AllowedPages = selpages
				}

				groupManagement.SaveGroup(global.GroupConfig.File)
			}
		}
	}
}

func deleteGroup(u string) {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	var gl []programTypes.Group
	n := make([]programTypes.Group, 1)

	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		if u != global.GroupConfig.File.Group[i].UUID {

			n[0] = global.GroupConfig.File.Group[i]

			gl = append(gl, n...)
		}
	}

	global.GroupConfig.File.Group = gl

	groupManagement.SaveGroup(global.GroupConfig.File)
}
