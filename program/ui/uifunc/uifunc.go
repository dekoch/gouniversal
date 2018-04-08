package uifunc

import (
	"gouniversal/program/global"
	"gouniversal/shared/functions"
	"gouniversal/shared/types"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LoginNameToUUID(user string) string {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if user == global.UserConfig.File.User[i].LoginName {

			return global.UserConfig.File.User[i].UUID
		}
	}

	return ""
}

func GetUserWithUUID(u string) types.User {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	for i := 0; i < len(global.UserConfig.File.User); i++ {

		if u == global.UserConfig.File.User[i].UUID {

			return global.UserConfig.File.User[i]
		}
	}

	var nu types.User

	return nu
}

func CheckLogin(user string, pwd string) bool {

	if functions.IsEmpty(user) == false &&
		functions.IsEmpty(pwd) == false {

		global.UserConfig.Mut.Lock()
		defer global.UserConfig.Mut.Unlock()

		for i := 0; i < len(global.UserConfig.File.User); i++ {

			if user == global.UserConfig.File.User[i].LoginName {

				return CheckPasswordHash(pwd, global.UserConfig.File.User[i].PWDHash)
			}
		}
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RemovePFromPath(path string) string {

	index := strings.Index(path, "$")

	var p string
	if index > 0 {
		p = path[:index]
	} else {
		p = path
	}

	return p
}
