package html

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/amery/file2go/static"
)

// FuncMap
func (c Collection) Funcs(funcs template.FuncMap) *template.Template {
	return c.root.Funcs(funcs)
}

func (c Collection) BindStaticCollection(hashify bool, sc static.Collection) *template.Template {
	funcMap := sc.NewFuncMap(hashify)
	return c.root.Funcs(template.FuncMap(funcMap))
}

// Parse
func (c Collection) Parse() error {

	for name, f := range c.files {

		// decode
		buf := new(strings.Builder)
		r, err := f.NewReader2("gzip")
		if err != nil {
			return err
		}

		_, err = io.Copy(buf, r)
		if err != nil {
			return err
		}

		// compile, and bind
		_, err = c.root.New(name).Parse(buf.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// Access
func (c Collection) Template(name string) (*template.Template, error) {
	if t := c.root.Lookup(name); t != nil {
		return t, nil
	} else {
		err := fmt.Errorf("html/template: %q not found")
		return nil, err
	}
}
