package jsonapi_test

import "strings"

func makeOneLine(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "\t", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)

	for {
		str2 := strings.Replace(str, "  ", " ", -1)
		if str == str2 {
			return str
		}
		str = str2
	}
}

func makeOneLineNoSpaces(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return strings.Replace(str, " ", "", -1)
}
