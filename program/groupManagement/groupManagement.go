package groupManagement

import (
	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
)

func IsPageAllowed(path string, gid string, checkState bool) bool {

	groups := global.GroupConfig.List()

	for i := 0; i < len(groups); i++ {

		g := groups[i]

		if gid == g.UUID {

			if checkState {
				// if group is not active
				if g.State != 1 {
					return false
				}
			}

			for p := 0; p < len(g.AllowedPages); p++ {

				allowed := uifunc.RemovePFromPath(g.AllowedPages[p])
				requested := uifunc.RemovePFromPath(path)

				if requested == allowed {

					return true
				}
			}

			return false
		}
	}

	return false
}
