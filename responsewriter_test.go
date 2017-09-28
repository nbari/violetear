package violetear

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriterStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	expect(t, rw.Status(), 0)

	rw.Write([]byte(""))
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 0)
}

func TestResponseWriterSize(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	rw.Write([]byte("日本語"))
	expect(t, rw.Size(), 9)

	rw.Write([]byte("a"))
	expect(t, rw.Size(), 10)
}

func TestResponseWriterHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	expect(t, len(rec.Header()), len(rw.Header()))
}

func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	rw.Write([]byte("Hello world"))
	rw.Write([]byte(". !"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello world. !")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 14)
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rw.Status(), 404)
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}
