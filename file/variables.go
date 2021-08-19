package file

import (
	"unicode"
)

func Varify(public bool, fname string) string {

	var capital, first bool

	buf := (make([]rune, 0, len(fname)+1))

	first = true

	for _, c := range fname {

		switch c {
		case '.', '/', '_', '-', ' ':
			capital = true
		default:

			if first {
				first = false // just once

				if !unicode.IsLetter(c) {
					// protect against invalid first runes
					var r rune
					if public {
						r = 'N'
					} else {
						r = 'n'
					}

					buf = append(buf, r)
				} else if public {
					c = unicode.ToUpper(c)
				} else {
					c = unicode.ToLower(c)
				}

			} else if capital {
				c = unicode.ToUpper(c)
			} else {
				c = unicode.ToLower(c)
			}

			capital = false
			buf = append(buf, c)
		}
	}

	return string(buf)
}
