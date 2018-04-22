package userManagement

import (
	"gouniversal/program/global"
	"gouniversal/program/groupManagement"
	"gouniversal/program/userConfig"
)

func SelectUser(uid string) userConfig.User {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	for u := 0; u < len(global.UserConfig.File.User); u++ {

		// search user with UUID
		if uid == global.UserConfig.File.User[u].UUID {

			return global.UserConfig.File.User[u]
		}
	}

	var user userConfig.User
	user.State = -1
	return user
}

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
