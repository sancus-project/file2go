package html

import (
	"html/template"

	"go.sancus.dev/file2go/file"
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

func (c *Collection) Clone() (c2 *Collection, err error) {
	var r2 *template.Template

	if r2, err = c.root.Clone(); err == nil {

		c2 = &Collection{
			root:  r2,
			files: c.files,
		}
	}

	return
}

// Extend() adds elements from one collection that aren't present already
func (c *Collection) Extend(b *Collection) {
	for k, v := range b.files {
		if _, ok := c.files[k]; !ok {
			c.files[k] = v
		}
	}
}

// Add() adds elements from one collection, and replaces those using the same key
func (c *Collection) Add(b *Collection) {
	for k, v := range b.files {
		c.files[k] = v
	}
}
