package violetear

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestResponseWriterLogger499(t *testing.T) {
	router := New()
	router.Verbose = false
	router.LogRequests = true
	router.Logger = func(w *ResponseWriter, r *http.Request) {
		expect(t, w.Status(), 499)
	}
	router.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	})
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := &http.Client{
		Timeout: time.Duration(time.Millisecond),
	}
	client.Get(ts.URL)
}

func TestResponseWriterXX(t *testing.T) {
	tt := []struct {
		name          string
		path          string
		reqMethod     string
		handlerMethod string
		rid           string
		ridValue      string
		code          int
		logRequests   bool
		logger        bool
	}{
		{"no logger", "/test", "GET", "GET", "rid", "123", 200, false, false},
		{"no logger 405", "/test", "GET", "POST", "rid", "123", 405, false, false},
		{"logger", "/test", "GET", "GET", "rid", "123", 200, true, true},
		{"logger 405", "/test", "GET", "POST", "rid", "123", 405, true, true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := New()
			if tc.logger {
				router.Logger = func(w *ResponseWriter, r *http.Request) {
					expect(t, r.URL.String(), tc.path)
					expect(t, w.RequestID(), tc.ridValue)
					expect(t, w.Status(), tc.code)
				}
			}
			router.RequestID = tc.rid
			router.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				expect(t, w.Header().Get(tc.rid), tc.ridValue)
			}, tc.handlerMethod)
			router.LogRequests = tc.logRequests
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.reqMethod, tc.path, nil)
			req.Header.Set(tc.rid, tc.ridValue)
			router.ServeHTTP(w, req)
			res := w.Result()
			expect(t, res.StatusCode, tc.code)
		})
	}
}
