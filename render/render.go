package render

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/amery/file2go/file"
	"github.com/amery/file2go/render/static"
)

type Renderer interface {
	Render(fout *os.File, varname string) error
	AddFile(fname string) error
}

func (c Config) Render(files []string) (err error) {
	var s []string
	var fname, varname, mode, pkg string
	var f *os.File
	var r Renderer

	// package
	pkg = c.Package
	if len(pkg) == 0 {
		return fmt.Errorf("Package name missing")
	}

	// mode
	switch c.Template {
	case "static", "none", "":
		r, err = static.NewStaticRenderer(files)
	default:
		return fmt.Errorf("Invalid Template mode %q", c.Template)
	}

	// output
	fname = c.Output
	if fname == "-" {
		fname = ""
	}

	// collection name
	varname = c.Varname
	if len(varname) > 0 {
		// take what's given
	} else if fname == "" {
		// default
		varname = "Files"
	} else {
		// take from filename
		ext := filepath.Ext(fname)
		varname = filepath.Base(fname)
		varname = file.Varify(true, varname[:len(varname)-len(ext)])
	}

	// create output
	fname = c.Output
	if len(fname) > 0 && fname != "-" {
		// temporary output, on the same directory for atomicity
		dir, fn := filepath.Split(fname)
		if dir == "" {
			dir = "./"
		}

		// dir/foo.go -> dir/.foo.go~
		tmpname := fmt.Sprintf("%s.%s~", dir, fn)

		f, err = os.OpenFile(tmpname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return
		}
		defer f.Close()
	} else {
		f = os.Stdout
		fname = ""
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

	// content
	if err = r.Render(f, varname); err != nil {
		// failed to render
		return
	} else if fname == "" {
		// fine, but no actual file
		return
	} else if err = f.Sync(); err != nil {
		// failed to flush
		return
	} else if err = os.Rename(f.Name(), fname); err != nil {
		// failed to rename
		return
	} else {
		// done.
		return
	}
}
