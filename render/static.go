package render

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/amery/file2go/static"
)

type StaticRendererFile struct {
	Name      string
	Hashified string
	Varname   string

	Content *static.Content
}

type StaticRenderer struct {
	Files map[string]*StaticRendererFile
	Names []string
}

func (r *StaticRenderer) Render(fout *os.File, files []string) error {

	// Initialise
	r.Files = make(map[string]*StaticRendererFile, len(files))
	r.Names = make([]string, 0, len(files))

	// Load files
	for _, fn := range files {
		if err := r.AddFile(fn); err != nil {
			return err
		}
	}

	r.Hashify()
	sort.Strings(r.Names)
	return r.Write(fout)
}

func (r *StaticRenderer) AddContent(fname string, blob *static.Content) error {
	var varname string

	// variable
	varname = "f_" + fname
	for _, c := range []string{".", "/", "-", " "} {
		varname = strings.Replace(varname, c, "_", -1)
	}

	o := &StaticRendererFile{
		Name:    fname,
		Varname: varname,
		Content: blob,
	}

	r.Files[fname] = o
	r.Names = append(r.Names, fname)
	return nil
}

func (r *StaticRenderer) AddFile(fname string) error {
	// input file
	if blob, err := static.NewContent(fname); err == nil {
		return r.AddContent(fname, blob)
	} else {
		return err
	}
}

func (r *StaticRenderer) Hashify() (err error) {
	files := make(map[string]*static.Content, len(r.Files))

	for fn0, v := range r.Files {
		files[fn0] = v.Content
	}

	m, _ := static.Hashify(files)
	for fn0, v := range r.Files {
		fn1 := m[fn0]
		if fn1 == fn0 {
			log.Printf("Hashify: %q", fn0)
		} else {
			log.Printf("Hashify: %q -> %q", fn0, fn1)
		}

		v.Hashified = fn1
	}

	return
}

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
		o := r.Files[fn0]
		v := o.Content

		if _, err := fmt.Fprintf(fout, "\n// %s\nvar %s = ", fn0, o.Varname); err != nil {
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
		_, err = fmt.Fprintf(fout, "\t%s[%q] = &%s\n", name, "/"+fn0, v)

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

		_, err = fmt.Fprintf(fout, "\t%s[%q] = &%s\n", name, "/"+fn1, v)
	}

	return
}

func (r *StaticRenderer) WriteEpilogue(fout *os.File) (err error) {

	_, err = fout.WriteString(`
var Files map[string]*static.Content
var HashifiedFiles map[string]*static.Content

func Handler(hashify bool, next http.Handler) http.Handler {
	var files map[string]*static.Content
	if hashify {
		files = HashifiedFiles
	} else {
		files = Files
	}

	return static.Handler(files, next)
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

	_, err = fout.WriteString("}")
	return
}
