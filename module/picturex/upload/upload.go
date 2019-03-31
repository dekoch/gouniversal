package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dekoch/gouniversal/module/picturex/global"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/google/uuid"
)

func LoadConfig() {

	http.HandleFunc("/picturex/upload/", handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	uid := r.FormValue("uuid")
	token := r.FormValue("token")
	pair := r.FormValue("pair")

	if global.Tokens.Check(uid, token) == false {
		fmt.Fprintf(w, "wrong token")
		return
	}

	path := global.Config.RawFileRoot

	err := functions.CreateDir(path)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	u := uuid.Must(uuid.NewRandom())
	name := u.String() + header.Filename

	out, err := os.Create(path + name)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	global.Tokens.Remove(uid)

	global.PairList.SetPicture(pair, uid, name)

	fmt.Fprintf(w, "<html><head><meta http-equiv=\"refresh\" content=\"0; url=/app\" /></head></html>")
}
