package errors

import (
	"fmt"
	"net/http"
)

var (
	ErrNotFound = &HandlerError{Code: http.StatusNotFound}
)

func ErrorText(code int) string {
	text := http.StatusText(code)

	if len(text) == 0 {
		text = fmt.Sprintf("Unknown Error %d", code)
	} else if code >= 400 {
		text = fmt.Sprintf("%s (Error %d)", text, code)
	}

	return text
}
