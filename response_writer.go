package violetear

import "net/http"

// LogResponseWritter wraps the standard http.ResponseWritter allowing for more
// verbose logging
type ResponseWritter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWritter(w http.ResponseWriter) *ResponseWritter {
	return &ResponseWritter{w, 0, 0}
}

// Status provides an easy way to retrieve the status code
func (w *ResponseWritter) Status() int {
	return w.status
}

// Size provides an easy way to retrieve the response size in bytes
func (w *ResponseWritter) Size() int {
	return w.size
}

// Header returns & satisfies the http.ResponseWriter interface
func (w *ResponseWritter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write satisfies the http.ResponseWriter interface and
// captures data written, in bytes
func (w *ResponseWritter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	size, err := w.ResponseWriter.Write(data)
	w.size += size

	return size, err
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows us to cach the status code
func (w *ResponseWritter) WriteHeader(statusCode int) {

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
