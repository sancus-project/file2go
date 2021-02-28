package html

import (
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/amery/file2go/file"
)

type Template struct {
	Name string
	Blob *file.Blob
}

func NewTemplate(name string, blob *file.Blob) Template {

	return Template{
		Name: name,
		Blob: blob,
	}
}

type Collection struct {
	tmpl map[string]*template.Template
}

func NewCollection2(entries ...Template) (*Collection, error) {
	root := template.New("")

	c := &Collection{
		tmpl: make(map[string]*template.Template, len(entries)),
	}

	for _, o := range entries {
		buf := new(strings.Builder)

		// read decoded into buffer
		r, err := o.Blob.NewReader2("gzip")
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(buf, r)
		if err != nil {
			return nil, err
		}

		// compile
		t, err := root.New(o.Name).Parse(buf.String())
		if err != nil {
			return nil, err
		}

		// and store
		c.tmpl[o.Name] = t
	}

	return c, nil
}

func NewCollection(entries ...Template) *Collection {

	c, err := NewCollection2(entries...)

	if err != nil {
		log.Panicln(err)
	}
	return c
}
