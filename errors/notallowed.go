package errors

import (
	"net/http"
	"strings"

	"go.sancus.dev/file2go/types"
)

type ErrNotAllowed struct {
	Method  string
	Allowed []string
}

func NotAllowed(method string, allowed ...string) types.Error {
	err := &ErrNotAllowed{
		Method:  method,
		Allowed: allowed,
	}
	return err
}

func (err *ErrNotAllowed) Status() int {
	if err.Method == "OPTIONS" {
		return http.StatusOK
	} else {
		return http.StatusMethodNotAllowed
	}
}

func (err *ErrNotAllowed) Error() string {
	return ErrorText(err.Status())
}

func (err *ErrNotAllowed) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	methods := append(err.Allowed, "OPTIONS")
	w.Header().Set("Allow", strings.Join(methods, ", "))

	if err.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		code := err.Status()
		http.Error(w, ErrorText(code), code)
	}
}
