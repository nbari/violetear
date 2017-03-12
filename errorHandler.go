package violetear

import "net/http"

type violetearError interface {
	error
	Status() int
}

// Error represents an error with an associated HTTP status code.
type Error struct {
	Code int
	Err  error
}

// Error return error message
func (e Error) Error() string {
	return e.Err.Error()
}

// Status return  HTTP status code.
func (e Error) Status() int {
	return e.Code
}

// ErrorHandler struct that returns error
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP allows ErrorHandler type to satisfy http.Handler.
func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		switch e := err.(type) {
		case violetearError:
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}
}
