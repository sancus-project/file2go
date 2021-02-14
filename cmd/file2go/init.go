package main

import (
	"fmt"
	"os"
	"text/template"
)

type file struct {
	Filename, Variable string
}

type data struct {
	Entries []file
}

func (c config) renderInit(fout *os.File, names, vars []string) error {
	d := data{
		Entries: make([]file, len(vars)),
	}

	for i, fname := range names {
		d.Entries[i] = file{
			Filename: fmt.Sprintf("/%s", fname),
			Variable: vars[i],
		}
	}

	t := template.Must(template.New("Init").Parse(`
var Files map[string]static.Content

func Handler(_ bool, next http.Handler) http.Handler {
	return static.Handler(Files, next)
}

func init() {
	Files = make(map[string]static.Content, {{len .Entries}})
{{range .Entries}}
	Files[{{printf "%q" .Filename}}] = {{.Variable}}{{end}}
}
`))
	return t.Execute(fout, d)
}
