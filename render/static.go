package render

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/amery/file2go/static"
)

type StaticRenderer struct{}

// CreateOutput
func (_ StaticRenderer) CreateOutput(fname string, pkg string, files []string) (f *os.File, err error) {
	if f, err = CreateOutput(fname, "", pkg, files); err != nil {
		return
	}

	// header
	if _, err = f.WriteString(`
import (
	"net/http"

	"github.com/amery/file2go/static"
)
`); err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}

// AddFile
func (_ StaticRenderer) AddFile(fout *os.File, fname string) (varname string, err error) {

	// variable
	varname = "f_" + fname
	for _, c := range []string{".", "/", "-", " "} {
		varname = strings.Replace(varname, c, "_", -1)
	}

	// input file
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer f.Close()

	// Template
	if _, err = fmt.Fprintf(fout, "\n// %s\nvar %s = ", fname, varname); err != nil {
		return
	}

	var content_type string

	if strings.HasSuffix(fname, ".css") {
		content_type = "text/css; charset=utf-8"
	} else if strings.HasSuffix(fname, ".js") {
		content_type = "application/javascript; charset=utf-8"
	}
	err = static.WriteContent(fout, f, content_type)

	if err == nil {
		_, err = fout.WriteString("\n")
	}
	return
}

// AddInit
type file struct {
	Filename, Variable string
}

type data struct {
	Entries []file
}

func (_ StaticRenderer) AddInit(fout *os.File, files []string, vars []string) error {
	d := data{
		Entries: make([]file, len(vars)),
	}

	for i, fname := range files {
		d.Entries[i] = file{
			Filename: fmt.Sprintf("/%s", fname),
			Variable: vars[i],
		}
	}

	t := template.Must(template.New("Init").Parse(`
var Files map[string]*static.Content

var hashifiedMap map[string]string
var hashifiedFiles map[string]*static.Content

func HashifiedMap(files map[string]*static.Content) map[string]string {
	if hashifiedMap == nil {
		hashifiedMap, hashifiedFiles = static.Hashify(files)
	}
	return hashifiedMap
}

func HashifiedFiles(files map[string]*static.Content) map[string]*static.Content {
	if hashifiedFiles == nil {
		hashifiedMap, hashifiedFiles = static.Hashify(files)
	}
	return hashifiedFiles
}

func Handler(hashify bool, next http.Handler) http.Handler {
	var files map[string]*static.Content
	if hashify {
		files = HashifiedFiles(Files)
	} else {
		files = Files
	}

	return static.Handler(files, next)
}

func init() {
	Files = make(map[string]*static.Content, {{len .Entries}})
{{range .Entries}}
	Files[{{printf "%q" .Filename}}] = &{{.Variable}}{{end}}
}
`))
	return t.Execute(fout, d)
}
