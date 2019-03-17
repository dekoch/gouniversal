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

	switch nav.GetNextPage() {
	case "User":
		user.Render(page, nav, r)

	case "Group":
		group.Render(page, nav, r)

	default:
		nav.RedirectPath("404", true)
	}
}
