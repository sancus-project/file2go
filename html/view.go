package html

import (
	"bytes"
	"html/template"
	"net/http"
)

// github.com/go-chi/render.Renderer
type View struct {
	tmpl *template.Template
	data interface{}
}

func (t View) Render(w http.ResponseWriter, _ *http.Request) error {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, t.data); err != nil {
		return err
	}
	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

// html/template
func (t View) Name() string {
	return t.tmpl.Name()
}

func (t View) Template() *template.Template {
	return t.tmpl
}

// Access
func (c Collection) View(name string, data interface{}) (v View, err error) {
	var t *template.Template

	t, err = c.Template(name)

	if err == nil {
		v = View{
			tmpl: t,
			data: data,
		}
	}

	return
}
