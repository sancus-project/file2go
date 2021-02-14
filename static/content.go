package static

import (
	"fmt"
	"net/http"
)

type Content struct {
	Body         Blob
	ContentType  string
	CacheControl string
	MaxAge       int
	Sha1sum      string
}

func (o Content) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	b := o.Body.RawBytes()

	h.Set("Content-Length", fmt.Sprintf("%v", len(b)))
	h.Set("Content-Type", o.ContentType)
	h.Set("Content-Encoding", o.Body.Encoding())

	if len(o.CacheControl) > 0 {
		h.Set("Cache-Control", o.CacheControl)
	} else {
		max_age := 24 * 60 * 60

		if o.MaxAge > 0 {
			max_age = o.MaxAge
		}

		h.Set("Cache-Control", fmt.Sprintf("max-age=%v", max_age))
	}
	if len(o.Sha1sum) > 0 {
		h.Set("ETag", o.Sha1sum)
	}

	w.Write(b)
}
