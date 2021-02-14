package render

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/amery/file2go/config"
)

type Renderer interface {
	CreateOutput(name string, pkg string, files []string) (*os.File, error)
	AddFile(fout *os.File, name string) (string, error)
	AddInit(fout *os.File, files []string, vars []string) error
}

func RenderConfig(c config.Render, files []string) {
	var fout *os.File
	var r Renderer
	var err error

	switch c.Template {
	default:
		r = StaticRenderer{}
	}

	sort.Strings(files)
	vars := make([]string, len(files))

	if len(c.Output) > 0 {
		// single output
		fout, err = r.CreateOutput(c.Output, c.Package, files)
		if err != nil {
			panic(err)
		}
	}

	for i, fname := range files {
		if len(c.Output) == 0 {
			// multi output
			if fout != nil {
				fout.Close()
			}

			fout, err = r.CreateOutput(fname+".go", c.Package, []string{fname})
			if err != nil {
				panic(err)
			}
		}

		v, err := r.AddFile(fout, fname)
		if err != nil {
			panic(err)
		} else {
			vars[i] = v
		}
	}

	if len(c.Output) > 0 {
		if err = r.AddInit(fout, files, vars); err != nil {
			panic(err)
		}
	}

	if fout != nil {
		fout.Close()
	}
}

func CreateOutput(fname string, mode string, pkg string, files []string) (f *os.File, err error) {
	var s []string

	if f, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return nil, err
	}

	s = append(s, fmt.Sprintf("//go:generate %s -p %s", os.Args[0], pkg))
	if len(fname) > 0 {
		s = append(s, fmt.Sprintf("-o %s", fname))
	}
	if len(mode) > 0 {
		s = append(s, fmt.Sprintf("-T %s", mode))
	}
	s = append(s, files...)

	if _, err = f.WriteString(strings.Join(s, " ")); err != nil {
		f.Close()
		return nil, err
	}

	if _, err = fmt.Fprintf(f, "\n\npackage %s\n", pkg); err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}
