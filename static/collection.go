package static

import (
	"net/http"
)

type Entry struct {
	Name      string
	Hashified string
	Content   *Content
}

func NewEntry(name, hashified string, content *Content) Entry {
	return Entry{
		Name:      name,
		Hashified: hashified,
		Content:   content,
	}
}

type Collection struct {
	Files     map[string]*Content
	Hashified map[string]*Content
	Redirects map[string]string
}

func NewCollection(entries ...Entry) *Collection {
	l := len(entries)

	c := &Collection{
		Files:     make(map[string]*Content, l),
		Hashified: make(map[string]*Content, l),
		Redirects: make(map[string]string, l),
	}

	for _, o := range entries {
		fn0 := o.Name
		fn1 := o.Hashified

		if len(fn1) == 0 {
			fn1 = fn0
		} else if fn0 != fn1 {
			c.Redirects[fn0] = fn1
		}

		c.Files[fn0] = o.Content
		c.Hashified[fn1] = o.Content
	}

	return c
}

// Extend() adds elements from one collection that aren't present already
func (c *Collection) Extend(b *Collection) {
	for fn0, fn1 := range b.Redirects {
		if o, ok := b.Files[fn0]; ok {

			if _, ok := c.Redirects[fn0]; !ok {
				// we didn't have a redirect for fn0, add it
				c.Redirects[fn0] = fn1
			}

			if _, ok := c.Files[fn0]; !ok {
				// we didn't have data for fn0, add it
				c.Files[fn0] = o
			}

			if _, ok := c.Hashified[fn1]; !ok {
				// we didn't have data for fn1, add it
				c.Hashified[fn1] = o
			}
		}
	}

	for fn0, o := range b.Files {
		if _, ok := c.Files[fn0]; !ok {
			// we didn't have data for fn0, add it
			c.Files[fn0] = o
		}
	}

	for fn1, o := range b.Hashified {
		if _, ok := c.Hashified[fn1]; !ok {
			// we didn't have data for fn1, add it
			c.Hashified[fn1] = o
		}
	}
}

// Add() adds elements from one collection, replacing those using the same key
func (c *Collection) Add(b *Collection) {
	for fn0, fn1 := range b.Redirects {
		c.Redirects[fn0] = fn1
	}
	for fn0, o := range b.Files {
		c.Files[fn0] = o
	}
	for fn1, o := range b.Hashified {
		c.Hashified[fn1] = o
	}
}

func (c *Collection) Handler(hashify bool, next http.Handler) http.Handler {
	var files map[string]*Content
	var redirects map[string]string

	if hashify {
		files = c.Hashified
		redirects = c.Redirects
	} else {
		files = c.Files
	}

	return Handler(files, redirects, next)
}

func (c *Collection) Middleware(hashify bool) func(http.Handler) http.Handler {
	var files map[string]*Content
	var redirects map[string]string

	if hashify {
		files = c.Hashified
		redirects = c.Redirects
	} else {
		files = c.Files
	}

	return func(next http.Handler) http.Handler {
		return Handler(files, redirects, next)
	}
}
