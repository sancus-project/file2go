package html

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"sort"

	"go.sancus.dev/file2go/file"
	tmpl "go.sancus.dev/file2go/template"
)

// Renderer
type HtmlRenderer struct {
	Files map[string]HtmlTemplateFile
	Names []string
}

func NewHtmlRenderer(files []string) (*HtmlRenderer, error) {
	r := &HtmlRenderer{}
	r.Files = make(map[string]HtmlTemplateFile, len(files))
	r.Names = make([]string, 0, len(files))

	for _, fn := range files {
		if err := r.AddFile(fn); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *HtmlRenderer) validate(b []byte) error {
	s := string(b)
	_, err := template.New("-").Parse(s)
	return err
}

func (r *HtmlRenderer) AddFile(fname string) error {
	encoder := &tmpl.TemplateEncoder{
		Validator: r.validate,
	}
	o := HtmlTemplateFile{}

	if err := o.Load(fname, encoder); err != nil {
		return err
	}

	o.Name = fname

	// remove extension
	ext := path.Ext(fname)
	key := fname[:len(fname)-len(ext)]

	o.Varname = file.Varify(false, key)

	r.Files[key] = o
	r.Names = append(r.Names, key)
	return nil
}

func (r *HtmlRenderer) Render(fout *os.File, varname string) error {

	sort.Strings(r.Names)

	if err := r.writePrologue(fout); err != nil {
		return err
	}
	if err := r.writeFiles(fout); err != nil {
		return err
	}
	if err := r.writeEpilogue(fout, varname); err != nil {
		return err
	}
	return nil
}

func (r *HtmlRenderer) writePrologue(fout *os.File) error {
	var file string

	if len(r.Names) > 0 {
		file = "\"go.sancus.dev/file2go/file\"\n\t"
	}

	_, err := fmt.Fprintf(fout, `
import (
	%s"go.sancus.dev/file2go/html"
)
`, file)
	return err
}
func (r *HtmlRenderer) writeFiles(fout *os.File) error {
	for _, key := range r.Names {
		v := r.Files[key]

		if _, err := fmt.Fprintf(fout, "\n// %s\nvar %s = ", v.Name, v.Varname); err != nil {
			return err
		}

		if err := v.Render(fout, "", 8); err != nil {
			return err
		}
	}

	return nil
}

func (r *HtmlRenderer) writeEpilogue(fout *os.File, varname string) error {

	var nl string

	if len(r.Names) > 0 {
		nl = "\n"
	}

	_, err := fmt.Fprintf(fout, "\nvar %s = html.NewCollection(%s", varname, nl)
	if err != nil {
		return err
	}

	for _, key := range r.Names {
		v := r.Files[key]

		_, err = fmt.Fprintf(fout, "\thtml.NewTemplate(%q, %s),\n", key, v.Varname)

		if err != nil {
			return err
		}
	}

	_, err = fout.WriteString(")\n")
	return err
}
