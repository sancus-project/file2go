package static

import (
	"fmt"
	"os"
)

type StaticRenderer struct {
	Files    map[string]*StaticRendererFile
	Names    []string
}

func NewStaticRenderer(files []string) (*StaticRenderer, error) {
	r := &StaticRenderer{}

	// Initialise
	r.Files = make(map[string]*StaticRendererFile, len(files))
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

		if _, err := fmt.Fprintf(fout, "\n// %s\nvar %s = ", fn0[1:], v.Varname()); err != nil {
			return err
		}

		if err := v.Render(fout, "", 8); err != nil {
			return err
		}
	}

	return nil
}

func (r *StaticRenderer) WriteEpilogue(fout *os.File) (err error) {

	_, err = fout.WriteString("\nvar Files = static.NewCollection(\n")
	if err != nil {
		return err
	}

	for _, fn0 := range r.Names {
		v := r.Files[fn0]
		fn1 := v.Hashified

		_, err = fmt.Fprintf(fout, "\tstatic.NewEntry(%q, %q, &%s),\n",
			fn0, fn1, v.Varname())
		if err != nil {
			return err
		}
	}

	_, err = fout.WriteString(`)

func Handler(hashify bool, next http.Handler) http.Handler {
	return Files.Handler(hashify, next)
}
`)
	if err != nil {
		return err
	}

	return nil
}
