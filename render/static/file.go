package static

import (
	"fmt"
	"os"

	"github.com/amery/file2go/file"
)

// file.File
type StaticRendererFile struct {
	file.File

	Hashified string
	Sha1sum   string
}

func (v StaticRendererFile) Render(fout *os.File, indent string, columns uint) (err error) {
	// Prologue
	_, err = fout.WriteString("static.Content{\n")
	if err != nil {
		return
	}

	// ContentType
	if len(v.ContentType) > 0 {
		_, err = fmt.Fprintf(fout, "%s\tContentType: %q,\n", indent, v.ContentType)
		if err != nil {
			return
		}
	}

	// Body
	if len(v.Body) > 0 {
		_, err = fmt.Fprintf(fout, "%s\tBody: []byte{\n", indent)
		if err != nil {
			return
		}

		last := len(v.Body) - 1
		for i, b := range v.Body {
			var pre, post string
			var col = uint(i) % columns

			if col != columns-1 {
				post = ""
			} else if i == last {
				post = ",\n\t" + indent
			} else {
				post = ",\n"
			}

			if col == 0 {
				// first column
				pre = "\t\t" + indent
			} else {
				pre = ", "
			}

			_, err = fmt.Fprintf(fout, "%s0x%02x%s", pre, b, post)
			if err != nil {
				return
			}
		}

		_, err = fout.WriteString("},\n")
		if err != nil {
			return
		}
	}

	// Sha1sum
	if len(v.Sha1sum) > 0 {
		_, err = fmt.Fprintf(fout, "%s\tSha1sum: %q,\n", indent, v.Sha1sum)

		if err != nil {
			return
		}
	}

	_, err = fout.WriteString("}\n")
	if err != nil {
		return
	}

	return
}

func (v StaticRendererFile) Varname() string {
	return file.Varify(v.Name)
}
