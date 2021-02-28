package html

import (
	"io"
	"html/template"
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

		// remove
		delete(c.files, name)

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

		// compile
		t, err := c.root.New(name).Parse(buf.String())
		if err != nil {
			return  err
		}

		// and store
		c.tmpl[name] = t
	}

	return nil
}

// Access
func (c Collection) Template(name string) *template.Template {
	if t, ok := c.tmpl[name]; ok {
		return t
	}
	return nil
}

// Execute
func (c Collection) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return c.root.ExecuteTemplate(wr, name, data)
}
