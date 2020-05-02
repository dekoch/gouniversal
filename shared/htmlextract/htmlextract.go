package htmlextract

import (
	"errors"
	"strings"
)

func Extract(s, prefix, suffix string) (string, error) {

	if strings.Index(s, prefix) == -1 {
		return s, errors.New("prefix not present")
	}

	if strings.Index(s, suffix) == -1 {
		return s, errors.New("suffix not present")
	}

	splitPrefix := strings.SplitAfter(s, prefix)

	for i := range splitPrefix {

		if strings.Index(splitPrefix[i], suffix) >= 0 {

			splitSuffix := strings.Split(splitPrefix[i], suffix)

			if len(splitSuffix) > 0 {
				return splitSuffix[0], nil
			}
		}
	}

	return s, errors.New("unknown error")
}
