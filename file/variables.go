package file

import (
	"unicode"
)

func Varify(public bool, fname string) string {

	var capital, first bool
	buf := (make([]rune, 0, len(fname)+1))

	if public {
		capital = true
	} else {
		first = true
	}

	for _, c := range fname {

		switch c {
		case '.', '/', '_', '-', ' ':
			capital = true
		default:

			if first || !capital {
				c = unicode.ToLower(c)
				first = false
			} else {
				c = unicode.ToUpper(c)
			}

			capital = false
			buf = append(buf, c)
		}
	}

	return string(buf)
}
