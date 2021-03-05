package html

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// FuncMap
func (c Collection) Funcs(funcs template.FuncMap) Collection {
	c.root.Funcs(funcs)
	return c
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
