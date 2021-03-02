package html

import (
	"html/template"

	"github.com/amery/file2go/file"
)

// Temporary wrapper
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
	root  *template.Template
	files map[string]*file.Blob
}

func NewCollection(entries ...Template) *Collection {

	c := &Collection{
		root:  template.New(""),
		files: make(map[string]*file.Blob, len(entries)),
	}

	// postpone compiling templates so we have time to add FuncMap
	for _, o := range entries {
		c.files[o.Name] = o.Blob
	}

	return c
}
