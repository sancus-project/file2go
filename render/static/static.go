package static

import (
	"fmt"
	"os"

	"github.com/amery/file2go/file"
)

type StaticRenderer struct {
	Files    map[string]*StaticRendererFile
	Redirect map[string]string
	Names    []string
}

func NewStaticRenderer(files []string) (*StaticRenderer, error) {
	r := &StaticRenderer{}

	// Initialise
	r.Files = make(map[string]*StaticRendererFile, len(files))
	r.Redirect = make(map[string]string, len(files))
	r.Names = make([]string, 0, len(files))

	// Load files
	for _, fn := range files {
		if err := r.AddFile(fn); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *StaticRenderer) Render(fout *os.File) error {
	r.Hashify()
	return r.Write(fout)
}

// AddFile
func (r *StaticRenderer) AddFile(fname string) error {
	encoder := &StaticRenderEncoder{}
	o := &StaticRendererFile{}

	// input file
	if err := o.Load(fname, encoder); err != nil {
		return err
	} else {
		fname = "/" + fname

		o.Name = fname
		o.Varname = file.Varify(fname)
		o.Sha1sum = encoder.Sha1sum()

		r.Files[fname] = o
		r.Names = append(r.Names, fname)
		return nil
	}
}

// Write
func (r *StaticRenderer) Write(fout *os.File) error {
	if err := r.WritePrologue(fout); err != nil {
		return err
	}
	if err := r.WriteFiles(fout); err != nil {
		return err
	}
	if err := r.WriteEpilogue(fout); err != nil {
		return err
	}
	return nil
}

func (r *StaticRenderer) WritePrologue(f *os.File) error {
	_, err := f.WriteString(`
import (
	"net/http"

	"github.com/amery/file2go/static"
)
`)

	return err
}

func (r *StaticRenderer) WriteFiles(fout *os.File) error {
	for _, fn0 := range r.Names {
		v := r.Files[fn0]

		if _, err := fmt.Fprintf(fout, "\n// %s\nvar %s = ", fn0[1:], v.Varname); err != nil {
			return err
		}

		if err := v.Render(fout, "", 8); err != nil {
			return err
		}
	}

	return nil
}

func (r *StaticRenderer) writeFilesInitTable(fout *os.File, name string) (err error) {
	_, err = fmt.Fprintf(fout, `
	// %s
	%s = make(map[string]*static.Content, %v)
`, name, name, len(r.Files))

	if err != nil {
		return
	}

	for _, fn0 := range r.Names {
		o := r.Files[fn0]
		v := o.Varname
		_, err = fmt.Fprintf(fout, "\t%s[%q] = &%s\n", name, fn0, v)

		if err != nil {
			return
		}
	}

	return
}

func (r *StaticRenderer) writeHashifiedInitTable(fout *os.File, name string) (err error) {
	_, err = fmt.Fprintf(fout, `
	// %s
	%s = make(map[string]*static.Content, %v)
`, name, name, len(r.Files))

	if err != nil {
		return
	}

	for _, fn0 := range r.Names {
		o := r.Files[fn0]
		v := o.Varname
		fn1 := o.Hashified

		_, err = fmt.Fprintf(fout, "\t%s[%q] = &%s\n", name, fn1, v)
	}

	return
}

func (r *StaticRenderer) writeRedirectInitTable(fout *os.File, name string) (err error) {
	_, err = fmt.Fprintf(fout, `
	// %s
	%s = make(map[string]string, %v)
`, name, name, len(r.Files))

	if err != nil {
		return
	}

	for _, fn0 := range r.Names {
		if fn1, ok := r.Redirect[fn0]; ok {
			_, err = fmt.Fprintf(fout, "\t%s[%q] = %q\n", name, fn0, fn1)
		}
	}

	return
}

func (r *StaticRenderer) WriteEpilogue(fout *os.File) (err error) {

	_, err = fout.WriteString(`
var Files map[string]*static.Content
var HashifiedFiles map[string]*static.Content
var Redirects map[string]string

func Handler(hashify bool, next http.Handler) http.Handler {
	var files map[string]*static.Content
	var redirects map[string]string
	if hashify {
		files = HashifiedFiles
		redirects = Redirects
	} else {
		files = Files
	}

	return static.Handler(files, redirects, next)
}

func init() {`)
	if err != nil {
		return
	}

	// Files
	if err = r.writeFilesInitTable(fout, "Files"); err != nil {
		return
	}

	// Hashified
	if err = r.writeHashifiedInitTable(fout, "HashifiedFiles"); err != nil {
		return
	}

	// Redirect
	if err = r.writeRedirectInitTable(fout, "Redirects"); err != nil {
		return
	}

	_, err = fout.WriteString("}\n")
	return
}
