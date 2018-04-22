package groupManagement

import (
	"gouniversal/program/global"
	"gouniversal/program/ui/uifunc"
)

func IsPageAllowed(path string, gid string, checkState bool) bool {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	for g := 0; g < len(global.GroupConfig.File.Group); g++ {

		if gid == global.GroupConfig.File.Group[g].UUID {

			if checkState {
				// if group is not active
				if global.GroupConfig.File.Group[g].State != 1 {
					return false
				}
			}

			for p := 0; p < len(global.GroupConfig.File.Group[g].AllowedPages); p++ {

				allowed := uifunc.RemovePFromPath(global.GroupConfig.File.Group[g].AllowedPages[p])
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
