package static

import (
	"fmt"
	"os"

	"go.sancus.dev/file2go/file"
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
	_, err = fmt.Fprintf(fout, "%s\tBody: ", indent)
	if err != nil {
		return
	}

	err = v.File.RenderBytes(fout, indent+"\t\t", columns)
	if err != nil {
		return
	}

	_, err = fout.WriteString(",\n")
	if err != nil {
		return
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
	return file.Varify(false, v.Name)
}
