package static

import (
	"net/http"
)

func serveFiles(w http.ResponseWriter, r *http.Request, files map[string]*Content, redirects map[string]string) bool {
	path := r.URL.Path

	// standarize path
	if len(path) == 0 {
		// empty path, skip
		return false
	} else if path[0] != '/' {
		path = "/" + path
	}

	if redirects != nil {
		if fn1, ok := redirects[path]; ok {
			http.Redirect(w, r, fn1, http.StatusTemporaryRedirect)
			// served
			return true
		}
	}

	if o, ok := files[path]; !ok {
		// unknown file, skip
		return false
	} else if r.Method == "GET" || r.Method == "HEAD" {
		o.ServeHTTP(w, r)
		// served
		return true
	} else {

		w.Header().Set("Allow", "OPTIONS, GET, HEAD")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		// handled
		return true
	}
}

func Handler(files map[string]*Content, redirects map[string]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !serveFiles(w, r, files, redirects) {
			next.ServeHTTP(w, r)
		}
	})
}
