package html

import (
	"fmt"
	"os"

	"github.com/amery/file2go/file"
)

type HtmlTemplateFile struct {
	file.File

	Varname string
}

func (f *HtmlTemplateFile) Render(fout *os.File, indent string, columns uint) error {
	var err error

	// Prologue
	_, err = fmt.Fprintf(fout, "&file.Blob{\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(fout, "%s\tBody: ", indent)
	if err != nil {
		return err
	}

	err = f.File.RenderBytes(fout, indent+"\t\t", columns)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(fout, ",\n%s}\n", indent)
	return err
}
