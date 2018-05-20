package settings

import (
	"net/http"

	"github.com/dekoch/gouniversal/program/ui/settings/group"
	"github.com/dekoch/gouniversal/program/ui/settings/user"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	user.RegisterPage(page, nav)
	group.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if nav.Path == "Program:Settings" {
		nav.NavigatePath("Program:Settings:General")
	}

	if nav.IsNext("User") {

		user.Render(page, nav, r)

	} else if nav.IsNext("Group") {

		group.Render(page, nav, r)
	} else {
		nav.RedirectPath("404", true)
	}
}
