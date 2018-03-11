package uifunc

import (
	"bytes"
	"gouniversal/program/global"
	"gouniversal/program/types"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func TemplToString(templ *template.Template, data interface{}) string {

	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {

	}

	return tpl.String()
}

type t int

const (
	STRING t = 1 + iota
	INT
)

func CheckInput(i string, typ t) bool {
	if typ == STRING {
		if len(strings.TrimSpace(i)) != 0 {
			return true
		}
	}

	if typ == INT {
		if len(strings.TrimSpace(i)) != 0 {
			return true
		}
	}

	return false
}

func CheckFormInput(key string, r *http.Request) string {

	input := r.FormValue(key)

	//if govalidator.IsAlphanumeric(input) {
	return input
	//}

	//return ""
}

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

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	if CheckInput(user, STRING) &&
		CheckInput(pwd, STRING) {

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
