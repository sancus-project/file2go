package static

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

type View struct {
	files     map[string]*Content
	redirects map[string]string
}

func (c Collection) View(hashify bool) (v *View) {

	if hashify {
		v = &View{
			files:     c.Hashified,
			redirects: c.Redirects,
		}
	} else {
		v = &View{
			files: c.Files,
		}
	}
	return
}

type PageInfo struct {
	prefix  string
	path    string
	content *Content
}

func (p PageInfo) Location() string {
	return p.prefix + p.path
}

func (p PageInfo) Canonical() string {
	return p.Location()
}

func (p PageInfo) Methods() []string {
	return []string{"HEAD", "GET"}
}

func (p PageInfo) MimeType() []string {
	return []string{p.content.ContentType}
}

func (p PageInfo) Handler() http.Handler {
	return p.content
}

type PageRedirect struct {
	page PageInfo
	path string
}

func (p PageRedirect) Location() string {
	return p.page.prefix + p.path
}

func (p PageRedirect) Canonical() string {
	return p.page.Canonical()
}

func (p PageRedirect) Handler() http.Handler {
	return errors.NewTemporaryRedirect(p.Canonical())
}

func (p PageRedirect) Error() error {
	return errors.NewTemporaryRedirect(p.Canonical())
}

// chi.Routes
func (v View) Routes() []chi.Route {
	var routes []chi.Route

	n := len(v.files) + len(v.redirects)

	if n > 0 {
		routes = make([]chi.Route, n)

		for k, o := range v.files {

			h := make(map[string]http.Handler, 2)
			h["GET"] = o
			h["HEAD"] = o

			r := chi.Route{
				Handlers: h,
				Pattern:  k,
			}

			routes = append(routes, r)
		}

		for k, loc := range v.redirects {
			o := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, loc, http.StatusTemporaryRedirect)
			})

			h := make(map[string]http.Handler, 2)
			h["GET"] = o
			h["HEAD"] = o

			r := chi.Route{
				Handlers: h,
				Pattern:  k,
			}

			routes = append(routes, r)
		}
	}

	return routes
}

func (v View) Middlewares() (m chi.Middlewares) {
	return
}

func (v View) Match(rctx *chi.Context, method, path string) bool {
	if method == "GET" || method == "HEAD" {
		if _, ok := v.redirects[path]; ok {
			return true
		}

		if _, ok := v.files[path]; ok {
			return true
		}
	}
	return false
}

// http.Handler
func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := handleFiles(w, r, v.files, v.redirects); err != nil {
		errors.HandleError(w, r, err)
	}
}

// web.Handler
func (v View) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return handleFiles(w, r, v.files, v.redirects)
}

// web.RouterPageInfo
func (v View) PageInfo(r *http.Request) (interface{}, bool) {
	var prefix string
	var path string

	if rctx := context.RouteContext(r.Context()); rctx != nil {
		prefix = rctx.RoutePrefix
		path = rctx.RoutePath
	} else {
		path = r.URL.Path
	}

	if dest, ok := v.redirects[path]; ok {
		// Redirect
		if c, ok := v.files[dest]; ok {
			p := &PageRedirect{
				path: path,
				page: PageInfo{
					prefix:  prefix,
					path:    dest,
					content: c,
				},
			}

			return p, true
		}
	}

	if c, ok := v.files[path]; ok {
		// Actual
		p := &PageInfo{
			prefix:  prefix,
			path:    path,
			content: c,
		}

		return p, true
	}

	return nil, false
}
