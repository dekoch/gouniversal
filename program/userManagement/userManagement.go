package userManagement

import (
	"github.com/dekoch/gouniversal/program/groupManagement"
	"github.com/dekoch/gouniversal/program/userConfig"
)

func IsUserInGroup(gid string, user userConfig.User) bool {

	for i := 0; i < len(user.Groups); i++ {

		// search group with UUID
		if gid == user.Groups[i] {

			return true
		}
	}

	return false
}

func IsPageAllowed(path string, user userConfig.User) bool {

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
		if groupManagement.IsPageAllowed(path, user.Groups[i], true) {
			return true
		}
	}

	return false
}
