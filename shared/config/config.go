package config

// FileHeader default config file header
type FileHeader struct {
	HeaderVersion  float32
	FileName       string
	ContentName    string
	ContentVersion float32
	Comment        string
}

// BuildHeader builds a default config file header
func BuildHeader(filename string, conname string, conver float32, comment string) FileHeader {
	var h FileHeader

	h.HeaderVersion = 1.0

	h.FileName = filename
	h.ContentName = conname
	h.ContentVersion = conver
	h.Comment = comment

	return h
}

// CheckHeader returns true, if conname and Header ContentName is identical
func CheckHeader(fh FileHeader, conname string) bool {
	if fh.ContentName != conname {
		return false
	}

	return true
}
