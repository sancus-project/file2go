package main

import (
	"fmt"
	"github.com/amery/file2go/static"
	"os"
	"strings"
)

func (c config) processFile(fout *os.File, fname string) (varname string, err error) {
	// variable
	varname = "f_" + fname
	for _, c := range []string{".", "/", "-", " "} {
		varname = strings.Replace(varname, c, "_", -1)
	}

	// input file
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer f.Close()

	// Template
	if _, err = fmt.Fprintf(fout, "\n// %s\nvar %s = ", fname, varname); err != nil {
		return
	}

	var content_type string

	if strings.HasSuffix(fname, ".css") {
		content_type = "text/css"
	} else if strings.HasSuffix(fname, ".js") {
		content_type = "application/javascript"
	}
	err = static.WriteContent(fout, f, content_type)

	if err == nil {
		_, err = fout.WriteString("\n")
	}
	return
}
