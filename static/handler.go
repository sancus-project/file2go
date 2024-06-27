package static

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.sancus.dev/web/errors"
)

func handleFiles(w http.ResponseWriter, r *http.Request, files map[string]*Content, redirects map[string]string) error {
	var path string

	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		path = rctx.RoutePath
	}

	if path == "" {
		path = r.URL.Path
	}

	// standarize path
	switch {
	case path == "":
		return errors.ErrNotFound
	case path[0] != '/':
		path = "/" + path
	}

	if redirects != nil {
		if fn1, ok := redirects[path]; ok {
			http.Redirect(w, r, fn1, http.StatusTemporaryRedirect)
			// served
			return nil
		}
	}

	if o, ok := files[path]; !ok {
		// unknown file, skip
		return errors.ErrNotFound
	} else if r.Method == "GET" || r.Method == "HEAD" {
		o.ServeHTTP(w, r)
		// served
		return nil
	} else {
		return errors.MethodNotAllowed(r.Method, "GET", "HEAD")
	}
}

func Handler(files map[string]*Content, redirects map[string]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handleFiles(w, r, files, redirects); err != nil {
			errors.HandleMiddlewareError(w, r, err, next)
		}
	})
}
