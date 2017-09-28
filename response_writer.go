package violetear

import (
	"net/http"
	"time"
)

// ResponseWriter wraps the standard http.ResponseWriter allowing for more
// verbose logging
type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

// NewResponseWriter returns ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 0, 0, time.Now()}
}

// Status provides an easy way to retrieve the status code
func (w *ResponseWriter) Status() int {
	return w.status
}

// Size provides an easy way to retrieve the response size in bytes
func (w *ResponseWriter) Size() int {
	return w.size
}

// Start retrieve the start time
func (w *ResponseWriter) Start() time.Time {
	return w.start
}

// Write satisfies the http.ResponseWriter interface and
// captures data written, in bytes
func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows us to cach the status code
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
