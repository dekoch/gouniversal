package navigation

import (
	"strings"

	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/program/usermanagement"
	"github.com/dekoch/gouniversal/shared/sitemap"
)

type Navigation struct {
	Path        string
	CurrentPath string
	LastPath    string
	Redirect    string
	User        userconfig.User
	Guest       bool
	GodMode     bool
	Sitemap     sitemap.Sitemap
}

func (nav *Navigation) GetNextPage() string {

	ret := strings.Replace(nav.Path, nav.CurrentPath, "", 1)
	// remove prefix ":" from :Module:Settings
	if strings.HasPrefix(ret, ":") {
		ret = ret[1:]
	}
	// remove all after ":" or "$"
	index := strings.Index(ret, ":")
	if index >= 0 {
		ret = ret[:index]
	} else {
		index = strings.Index(ret, "$")
		if index >= 0 {
			ret = ret[:index]
		}
	}
	// remember current path
	if len(nav.CurrentPath) == 0 {
		nav.CurrentPath = ret
	} else {
		nav.CurrentPath += ":" + ret
	}

	return ret
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

	if len(path) == 0 {
		return
	}

	nav.LastPath = nav.Path

	if overwrite {
		nav.Redirect = path
	} else {
		if len(nav.Redirect) == 0 {
			nav.Redirect = path
		}
	}
}

// Parameter returns value from page parameter
// $Parametername=<value>
func (nav *Navigation) Parameter(name string) string {

	name = "$" + name + "="

	startIndex := strings.Index(nav.Path, name)
	if startIndex < 0 {
		return ""
	}

	ret := nav.Path[startIndex:]
	ret = strings.Replace(ret, name, "", 1)

	endIndex := strings.Index(ret, "$")
	if endIndex >= 0 {
		ret = ret[:endIndex]
	}

	return ret
}
