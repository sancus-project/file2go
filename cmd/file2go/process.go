package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func (c config) createOutput(fname string, args []string) (f *os.File, err error) {
	flags := os.O_WRONLY
	if c.appendOutput {
		f, _ = os.OpenFile(fname, flags|os.O_APPEND, 0644)
	}

	if f == nil {
		// new output
		f, err = os.OpenFile(fname, flags|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}

		// go:generate
		if !c.omitGoGenerate {
			var s []string
			s = append(s, fmt.Sprintf("//go:generate %s -p %s", os.Args[0], c.Package))
			if len(c.Output) > 0 {
				s = append(s, fmt.Sprintf("-o %s", c.Output))
			}
			s = append(s, args...)

			if _, err = f.WriteString(strings.Join(s, " ")); err != nil {
				f.Close()
				return nil, err
			}

			if _, err = f.WriteString("\n\n"); err != nil {
				f.Close()
				return nil, err
			}
		}

		// header
		s := `package %s

import (
	"github.com/amery/file2go/static"
)
`
		if _, err = fmt.Fprintf(f, s, c.Package); err != nil {
			f.Close()
			return nil, err
		}
	}

	return f, nil
}

func (c config) Process(args []string) {
	var fout *os.File
	var err error

	sort.Strings(args)
	vars := make([]string, len(args))

	if len(c.Output) > 0 {
		fout, err = c.createOutput(c.Output, args)
		if err != nil {
			panic(err)
		}
	}

	for i, fname := range args {
		if len(c.Output) == 0 {
			if fout != nil {
				fout.Close()
			}

			fout, err = c.createOutput(fname+".go", []string{fname})
			if err != nil {
				panic(err)
			}
		}

		v, err := c.processFile(fout, fname)
		if err != nil {
			panic(err)
		} else {
			vars[i] = v
		}
	}

	if len(c.Output) > 0 {
		if err = c.renderInit(fout, args, vars); err != nil {
			panic(err)
		}
	}

	if fout != nil {
		fout.Close()
	}
}
