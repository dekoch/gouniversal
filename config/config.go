package config

type FileHeader struct {
	HeaderVersion  float32
	FileName       string
	ContentName    string
	ContentVersion float32
	Comment        string
}

func BuildHeader(filename string, conname string, conver float32, comment string) FileHeader {
	var h FileHeader

	h.HeaderVersion = 1.0

	h.FileName = filename
	h.ContentName = conname
	h.ContentVersion = conver
	h.Comment = comment

	return h
}

func CheckHeader(fh FileHeader, conname string) bool {
	if fh.ContentName != conname {
		return false
	}

	return true
}
