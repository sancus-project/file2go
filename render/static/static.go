package static

import (
	"fmt"
	"os"
)

type StaticRenderer struct {
	Files map[string]*StaticRendererFile
	Names []string
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

func (r *StaticRenderer) Render(fout *os.File, varname string) error {
	r.Hashify()
	return r.render(fout, varname)
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
func (r *StaticRenderer) render(fout *os.File, varname string) error {
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

func (r *StaticRenderer) writePrologue(f *os.File) error {
	_, err := f.WriteString(`
import (
	"github.com/amery/file2go/static"
)
`)

	return err
}

func (r *StaticRenderer) writeFiles(fout *os.File) error {
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

func (r *StaticRenderer) writeEpilogue(fout *os.File, varname string) (err error) {

	_, err = fmt.Fprintf(fout, "\nvar %s = static.NewCollection(\n", varname)
	if err != nil {
		return
	}

	for _, fn0 := range r.Names {
		v := r.Files[fn0]
		fn1 := v.Hashified

		_, err = fmt.Fprintf(fout, "\tstatic.NewEntry(%q, %q, &%s),\n",
			fn0, fn1, v.Varname())
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintf(fout, ")\n")
	return
}
