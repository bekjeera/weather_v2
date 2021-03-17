package handler

import "net/http"

// StatusError define extended error with Status
type StatusError interface {
	Status() int
	error
}

// NewStatusError is a helper function
func NewStatusError(status int, message string) StatusError {
	return httpError{status, message}
}

// httpError satisfies StatusError interface
type httpError struct {
	status  int
	message string
}

func (e httpError) Status() int {
	return e.status
}

func (e httpError) Error() string {
	return e.message
}

// ErrorHandler is like a middleware for each handler
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// it should implement ServeHTTP method to be a handler type
func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}
	status := http.StatusInternalServerError
	if serr, ok := err.(StatusError); ok {
		status = serr.Status()
	}
	http.Error(w, err.Error(), status)
}
