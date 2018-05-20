package uifunc

import (
	"gouniversal/program/global"
	"gouniversal/shared/functions"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LoginNameToUUID(name string) string {

	user := global.UserConfig.List()

	for i := 0; i < len(user); i++ {

		if name == user[i].LoginName {

			return user[i].UUID
		}
	}

	return ""
}

func CheckLogin(name string, pwd string) bool {

	if functions.IsEmpty(name) == false &&
		functions.IsEmpty(pwd) == false {

		u, err := global.UserConfig.GetWithName(name)
		if err != nil {
			return false
		}

		return CheckPasswordHash(pwd, u.PWDHash)
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
