package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dekoch/gouniversal/shared/functions"
)

// http://blog.ralch.com/tutorial/golang-working-with-zip/
// https://gist.github.com/svett/424e6784facc0ba907ae

func Zipit(source, target string, exclude []string) error {

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return err
	}

	err := functions.CreateDir(target)
	if err != nil {
		return err
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			for i := range exclude {

				if strings.HasPrefix(path, exclude[i]) {
					return nil
				}
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})
}

func Unzip(archive, target string) error {

	if _, err := os.Stat(archive); os.IsNotExist(err) {
		return err
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	err = functions.CreateDir(target)
	if err != nil {
		return err
	}

	for _, file := range reader.File {

		path := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {

			err = os.MkdirAll(path, file.Mode())
			if err != nil {
				return err
			}

			continue
		}

		fileReader, err := file.Open()
		if err != nil {

			if fileReader != nil {
				fileReader.Close()
			}

			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()

			if targetFile != nil {
				targetFile.Close()
			}

			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()

			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return nil
}
