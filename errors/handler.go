package errors

import (
	"fmt"
	"net/http"
)

type HandlerError struct {
	Code int
}

func (err HandlerError) Status() int {
	if err.Code == 0 {
		return http.StatusOK
	} else {
		return err.Code
	}
}

func (err HandlerError) Error() string {
	return ErrorText(err.Status())
}

func (err HandlerError) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if code := err.Status(); code == http.StatusOK {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(code)

		fmt.Fprintln(w, ErrorText(code))
	}
}
