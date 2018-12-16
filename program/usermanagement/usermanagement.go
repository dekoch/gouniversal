package usermanagement

import (
	"github.com/dekoch/gouniversal/program/groupmanagement"
	"github.com/dekoch/gouniversal/program/userconfig"
)

func IsUserInGroup(gid string, user userconfig.User) bool {

	for i := 0; i < len(user.Groups); i++ {

		// search group with UUID
		if gid == user.Groups[i] {

			return true
		}
	}

	return false
}

func IsPageAllowed(path string, user userconfig.User) bool {

	// always allowed pages
	if path == "Account:Login" ||
		path == "Account:Logout" {

		return true
	}

	if path == "Program:Home" &&
		user.State == 1 {

		return true
	}

	// if user state is not set
	if user.State < 0 {
		return false
	}

	// test each group
	for i := 0; i < len(user.Groups); i++ {
		if groupmanagement.IsPageAllowed(path, user.Groups[i], true) {
			return true
		}
	}

	return false
}
