package pageGroupEdit

import (
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/groupManagement"
	"gouniversal/program/lang"
	"gouniversal/program/types"
	"gouniversal/program/ui/navigation"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/ui/uiglobal"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func RegisterPage(page *uiglobal.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program:Settings:Group:Edit", page.Lang.Settings.Group.GroupEdit.Title)
}

func Render(page *uiglobal.Page, nav *navigation.Navigation, r *http.Request) {

	type groupedit struct {
		Lang  lang.SettingsGroupEdit
		Group types.Group
		Pages template.HTML
	}
	var ge groupedit

	ge.Lang = page.Lang.Settings.Group.GroupEdit

	button := r.FormValue("edit")

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

	// copy groups from array
	global.GroupConfig.Mut.Lock()
	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		if id == global.GroupConfig.File.Group[i].UUID {

			ge.Group = global.GroupConfig.File.Group[i]
		}
	}
	global.GroupConfig.Mut.Unlock()

	pagelist := ""
	sm := nav.Sitemap.PageList()

	for i := 0; i < len(sm); i++ {

		pagelist += "<tr>"
		pagelist += "<td>" + sm[i] + "</td>"
		pagelist += "<td><input type=\"checkbox\" name=\"selectedpages\" value=\"" + sm[i] + "\""

		if groupManagement.IsPageAllowed(sm[i], ge.Group.UUID) {

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

	page.Content += uifunc.TemplToString(templ, ge)
}

func newGroup() string {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())

	newgroup := make([]types.Group, 1)
	newgroup[0].UUID = u.String()
	newgroup[0].Name = u.String()

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

	var gl []types.Group
	n := make([]types.Group, 1)

	for i := 0; i < len(global.GroupConfig.File.Group); i++ {

		if u != global.GroupConfig.File.Group[i].UUID {

			n[0] = global.GroupConfig.File.Group[i]

			gl = append(gl, n...)
		}
	}

	global.GroupConfig.File.Group = gl

	groupManagement.SaveGroup(global.GroupConfig.File)
}
