package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/amery/file2go/render/static"
)

type Renderer interface {
	Render(fout *os.File) error
	AddFile(fname string) error
}

func (c Config) Render(files []string) (err error) {
	var s []string
	var fname, mode, pkg string
	var f *os.File
	var r Renderer

	// turn config.Render into variables
	fname = c.Output
	pkg = c.Package

	if len(pkg) == 0 {
		return fmt.Errorf("Package name missing")
	}

	switch c.Template {
	case "static", "none", "":
		r, err = static.NewStaticRenderer(files)
	default:
		return fmt.Errorf("Invalid Template mode %q", c.Template)
	}

	// Create output
	if len(fname) > 0 {
		f, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			return
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}

	// go:generate
	s = append(s, fmt.Sprintf("//go:generate %s -p %s", os.Args[0], pkg))
	if len(fname) > 0 {
		s = append(s, fmt.Sprintf("-o %s", fname))
	}
	if len(mode) > 0 {
		s = append(s, fmt.Sprintf("-T %s", mode))
	}
	s = append(s, files...)

	if _, err = f.WriteString(strings.Join(s, " ")); err != nil {
		return
	}

	// package
	if _, err = fmt.Fprintf(f, "\n\npackage %s\n", pkg); err != nil {
		return
	}

	return r.Render(f)
}
