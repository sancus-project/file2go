package static

import (
	"net/http"
)

func serveFiles(w http.ResponseWriter, r *http.Request, files map[string]Content) bool {
	path := r.URL.Path

	// standarize path
	if len(path) == 0 {
		// empty path, skip
		return false
	} else if path[0] != '/' {
		path = "/" + path
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

func Handler(files map[string]Content, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !serveFiles(w, r, files) {
			next.ServeHTTP(w, r)
		}
	})
}
