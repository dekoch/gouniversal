package navigation

import (
	"strings"

	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/program/usermanagement"
	"github.com/dekoch/gouniversal/shared/sitemap"
)

type Navigation struct {
	Path           string
	CurrentPath    string
	LastPath       string
	Redirect       string
	PathAfterLogin string
	Home           string
	User           userconfig.User
	Guest          bool
	GodMode        bool
	Sitemap        sitemap.Sitemap
}

func (nav Navigation) CanGoBack() bool {

	if strings.Count(nav.CurrentPath, ":") > 0 {
		return true
	}

	return false
}

func (nav *Navigation) GoBack() {

	if nav.CanGoBack() {

		index := strings.LastIndex(nav.Path, ":")

		cnt := len(nav.Path)

		if cnt > 0 {
			nav.Path = nav.Path[:index]
		}
	}
}

func (nav *Navigation) Navigate(page string) {

	if len(page) > 0 {

		if strings.HasSuffix(nav.Path, ":"+page) == false {

			nav.Path += ":" + page
		}
	}
}

func (nav *Navigation) IsNext(page string) bool {

	if len(page) > 0 {

		nextPage := nav.CurrentPath + ":" + page

		if strings.HasPrefix(nextPage, ":") {
			nextPage = nextPage[1:]
		}

		if strings.HasPrefix(nav.Path, nextPage) {

			nav.CurrentPath += ":" + page

			if strings.HasPrefix(nav.CurrentPath, ":") {
				nav.CurrentPath = nav.CurrentPath[1:]
			}

			//console.Output(nav.CurrentPath, "navigation.IsNext()")

			return true
		}
	}

	return false
}

func (nav *Navigation) NavigatePath(path string) {

	if len(path) > 0 {

		if usermanagement.IsPageAllowed(path, nav.User) ||
			nav.GodMode {

			nav.Path = path
		}

		// debug
		//nav.Path = path
	}
}

func (nav *Navigation) RedirectPath(path string, overwrite bool) {

	nav.LastPath = nav.Path

	if overwrite {

		if len(path) > 0 {
			nav.Redirect = path
		}
	} else {

		if len(nav.Redirect) == 0 && len(path) > 0 {

			nav.Redirect = path
		}
	}
}

func (nav *Navigation) AfterLogin(path string) {

	if len(path) > 0 {
		nav.PathAfterLogin = path
	}
}

func (nav *Navigation) NavigateHome() {
	nav.NavigatePath(nav.Home)
}

// Parameter returns value from page parameter
func (nav *Navigation) Parameter(name string) string {

	par := ""

	name += "="

	index := strings.LastIndex(nav.Path, name)
	if index < 0 {
		return ""
	}

	cnt := len(nav.Path)

	if cnt > 0 {
		par = nav.Path[index:]

		par = strings.Replace(par, name, "", 1)
	}

	return par
}
