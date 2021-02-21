package file

import (
	"unicode"
)

func Varify(fname string) string {

	buf := (make([]rune, 0, len(fname)+1))
	capital := false
	first := true

	for _, c := range fname {

		switch c {
		case '.', '/', '-', ' ':
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
