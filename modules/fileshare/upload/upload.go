package upload

import (
	"fmt"
	"gouniversal/modules/fileshare/global"
	"gouniversal/shared/functions"
	"io"
	"net/http"
	"os"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	uid := r.FormValue("uuid")
	token := r.FormValue("token")

	if global.Tokens.Check(uid, token) == false {
		fmt.Fprintf(w, "wrong token")
		return
	}

	path := r.FormValue("path")
	path = global.Config.File.FileRoot + path

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

	out, err := os.Create(path + header.Filename)
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

	fmt.Fprintf(w, "<html><head><meta http-equiv=\"refresh\" content=\"0; url=/app\" /></head></html>")
}

func LoadConfig() {

	http.HandleFunc("/fileshare/upload/", handleRequest)
}
