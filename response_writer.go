package violetear

import (
	"golang.org/x/net/context"
	"net/http"
)

// ResponseWriter wraps the standard http.ResponseWriter allowing for more
// verbose logging
type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
	ctx    context.Context
}

// NewResponseWriter returns ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 0, 0, context.TODO()}
}

// Status provides an easy way to retrieve the status code
func (w *ResponseWriter) Status() int {
	return w.status
}

// Size provides an easy way to retrieve the response size in bytes
func (w *ResponseWriter) Size() int {
	return w.size
}

// Header returns & satisfies the http.ResponseWriter interface
func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
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

func (w *ResponseWriter) Get(s string) interface{} {
	return w.ctx.Value(s)
}

func (w *ResponseWriter) Set(k string, v interface{}) {
	w.ctx = context.WithValue(w.ctx, k, v)
}

func (w *ResponseWriter) SetParam(n string, v string) {
	param := w.Get(n)
	if param != nil {
		s := []interface{}{param}
		s = append(s, v)
		w.Set(n, s)
	} else {
		w.Set(n, v)
	}
}
