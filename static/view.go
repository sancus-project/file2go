package static

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

// http.Handle
func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !serveFiles(w, r, v.files, v.redirects) {
		http.NotFound(w, r)
	}
}
