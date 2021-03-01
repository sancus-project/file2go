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

func NewCollection(entries ...Entry) Collection {
	l := len(entries)

	c := Collection{
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

func (c Collection) Handler(hashify bool, next http.Handler) http.Handler {
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

func (c Collection) Middleware(hashify bool) func(http.Handler) http.Handler {
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
