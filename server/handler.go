package server

import "net/http"
import (
	errs "errors"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

//stackTracer is an error containing stack trace
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HTTPError represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type HTTPError interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

//NewStatusError creates new StatusError
func NewStatusError(code int, err string) StatusError {
	return StatusError{code, errs.New(err)}
}

//Error allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

//Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		if err, ok := err.(stackTracer); ok {

			stackTrace := make([]string, len(err.StackTrace()))
			for i, f := range err.StackTrace() {
				stackTrace[i] = fmt.Sprintf("%+s", f)
			}
			fmt.Println(strings.Join(stackTrace, "\n"))
		}

		switch e := errors.Cause(err).(type) {
		case HTTPError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			WriteJSON(e.Status(), map[string]string{"error": e.Error()}, w)
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
