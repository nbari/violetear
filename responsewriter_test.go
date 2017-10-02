package violetear

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriterStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, "")

	expect(t, rw.Status(), 200)

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

func TestResponseWriterLogger(t *testing.T) {
	mylogger := func(w *ResponseWriter, r *http.Request) {
		expect(t, r.URL.String(), "/test")
		expect(t, w.RequestID(), "123")
		expect(t, w.Size(), 11)
		expect(t, w.Status(), 200)
	}
	router := New()
	router.LogRequests = true
	router.RequestID = "rid"
	router.Logger = mylogger
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		expect(t, w.Header().Get("rid"), "123")
		w.Write([]byte("hello world"))
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("rid", "123")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
	expect(t, w.HeaderMap.Get("rid"), "123")
}

func TestResponseWriterLoggerStatus200(t *testing.T) {
	mylogger := func(w *ResponseWriter, r *http.Request) {
		expect(t, r.URL.String(), "/test")
		expect(t, w.RequestID(), "123")
		expect(t, w.Size(), 0)
		expect(t, w.Status(), 200)
	}
	router := New()
	router.LogRequests = true
	router.RequestID = "rid"
	router.Logger = mylogger
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		expect(t, w.Header().Get("rid"), "123")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("rid", "123")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
	expect(t, w.HeaderMap.Get("rid"), "123")
}

func TestResponseWriterLoggerStatus405(t *testing.T) {
	mylogger := func(w *ResponseWriter, r *http.Request) {
		expect(t, r.URL.String(), "/test")
		expect(t, w.RequestID(), "123")
		expect(t, w.Status(), 405)
	}
	router := New()
	router.LogRequests = true
	router.RequestID = "rid"
	router.Logger = mylogger
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		expect(t, w.Header().Get("rid"), "123")
	}, "POST")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("rid", "123")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 405)
	expect(t, w.HeaderMap.Get("rid"), "123")
}

func TestResponseWriterNoLogger(t *testing.T) {
	router := New()
	router.LogRequests = false
	router.RequestID = "rid"
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		expect(t, w.Header().Get("rid"), "123")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("rid", "123")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
	expect(t, w.HeaderMap.Get("rid"), "123")
}

func TestResponseWriterNoLogger455(t *testing.T) {
	router := New()
	router.LogRequests = false
	router.RequestID = "rid"
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		expect(t, w.Header().Get("rid"), "123")
	}, "POST")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("rid", "123")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 405)
	expect(t, w.HeaderMap.Get("rid"), "123")
}
