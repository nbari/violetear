package violetear

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/nbari/violetear/middleware"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	if a != b {
		t.Fatalf("Expected: %v (type %v)  Got: %v (type %v)  in %s:%d", b, reflect.TypeOf(b), a, reflect.TypeOf(a), fn, line)
	}
}

func expectDeepEqual(t *testing.T, a interface{}, b interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Expected: %v (type %v)  Got: %v (type %v)  in %s:%d", b, reflect.TypeOf(b), a, reflect.TypeOf(a), fn, line)
	}
}

type testRouter struct {
	path     string
	methods  string
	requests []testRequests
}

type testRequests struct {
	request string
	method  string
	expect  int
}

type testDynamicRoutes struct {
	name  string
	regex string
}

var dynamicRoutes = []testDynamicRoutes{
	{":uuid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`},
	{":ip", `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`},
}

var routes = []testRouter{
	{"/", "", []testRequests{
		{"/", "GET", 200},
	}},
	{"*", "GET", []testRequests{
		{"/a", "GET", 200},
		{"/a", "HEAD", 405},
		{"/a", "POST", 405},
	}},
	{"/:uuid", "GET, HEAD", []testRequests{
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25", "GET", 200},
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25", "HEAD", 200},
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25", "POST", 405},
	}},
	{"/:uuid/1/", "PUT", []testRequests{
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25/1", "PUT", 200},
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25/2", "GET", 404},
		{"/3B96853C-EF0B-44BC-8820-A982A5756E25/not_found/44", "GET", 404},
		{"/D0ABD486-B05A-436B-BBD1-E320CDC87916/1", "PUT", 200},
	}},
	{"/root", "GET,HEAD", []testRequests{
		{"/root", "GET", 200},
		{"/root", "HEAD", 200},
		{"/root", "OPTIONS", 405},
		{"/root", "POST", 405},
		{"/root", "PUT", 405},
	}},
	{"/root/:ip/", "GET", []testRequests{
		{"/root/10.0.0.0", "GET", 200},
		{"/root/172.16.0.0", "GET", 200},
		{"/root/192.168.0.1", "GET", 200},
		{"/root/300.0.0.0", "GET", 404},
	}},
	{"/root/:ip/aaa/", "GET", []testRequests{}},
	{"/root/:ip/aaa/:uuid", "GET", []testRequests{}},
	{"/root/:uuid/", "PATCH", []testRequests{
		{"/root/3B96853C-EF0B-44BC-8820-A982A5756E25", "GET", 405},
		{"/root/3B96853C-EF0B-44BC-8820-A982A5756E25", "PATCH", 200},
	}},
	{"/root/:uuid/-/:uuid", "GET", []testRequests{
		{"/root/22314BF-4A90-46C8-948D-5507379BD0DD/-/4293C253-6C7E-4B01-90F2-18203FAB2AEC", "GET", 404},
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/-/4293C253-6C7E-4B01-90F2-18203FAB2AE", "GET", 404},
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/-/4293C253-6C7E-4B01-90F2-18203FAB2AEF", "GET", 200},
		{"/root/E22314BF-4A90-46C8-948D-5507379BD0DD/-/4293C253-6C7E-4B01-90F2-18203FAB2AEC", "GET", 200},
	}},
	{"/root/:uuid/:uuid", "", []testRequests{
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AE", "GET", 404},
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF", "GET", 200},
	}},
	{"/root/:uuid/:uuid/end", "GET", []testRequests{
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/end", "GET", 200},
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/end-not-found", "GET", 404},
	}},
	{"/toor/", "GET", []testRequests{
		{"/toor", "GET", 200},
	}},
	{"/toor/aaa", "GET", []testRequests{
		{"/toor/aaa", "GET", 200},
		{"/toor/abc", "GET", 404},
	}},
	{"/toor/*", "GET", []testRequests{
		{"/toor/abc", "GET", 200},
		{"/toor/epazote", "GET", 200},
		{"/toor/naranjas", "GET", 200},
	}},
	{"/toor/1/2", "GET", []testRequests{
		{"/toor/1/2", "GET", 200},
	}},
	{"/toor/1/*", "GET", []testRequests{
		{"/toor/1/catch-me", "GET", 200},
		{"/toor/1/catch-me/too", "GET", 200},
		{"/toor/1/catch-me/too/foo/bar", "GET", 200},
	}},
	{"/toor/1/2/3", "GET", []testRequests{
		{"/toor/1/2/3", "GET", 200},
	}},
	{"/not-found", "GET", []testRequests{
		{"/toor/1/2/3/4", "GET", 404},
		{"catch_me", "GET", 200},
	}},
	{"/root/:uuid/:uuid/:ip/catch-me", "GET", []testRequests{}},
	{"/root/:uuid/:uuid/:ip/catch-me/*", "GET", []testRequests{}},
	{"/root/:uuid/:uuid/:ip/dont-wcatch-me", "GET", []testRequests{}},
	{"/root/:uuid/:uuid/:ip/dont-wcatch-me", "GET", []testRequests{}},
	{"/root/:uuid/:uuid/:ip/", "GET", []testRequests{
		{"/root/122314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8", "GET", 200},
		{"/root/122314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8/catch-me", "GET", 200},
		{"/root/122314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8/catch-me/also", "GET", 200},
		{"/root/122314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8/catch-me/also/a/b/c", "GET", 200},
		{"/root/122314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8/dont-catch-me", "GET", 404},
		{"/root/A22314BF-4A90-46C8-948D-5507379BD0DD/4293C253-6C7E-4B01-90F2-18203FAB2AEF/8.8.8.8", "GET", 200},
	}},
	{"/violetear/:ip/:uuid", "GET", []testRequests{
		{"/violetear/", "GET", 404},
		{"/violetear/127.0.0.1/", "GET", 404},
		{"/violetear/127.0.0.1/A22314BF-4A90-46C8-948D-5507379BD0DD/", "GET", 200},
		{"/violetear/127.0.0.1/A22314BF-4A90-46C8-948D-5507379BD0DD/not-found", "GET", 404},
	}},
	{"/:ip", "GET", []testRequests{
		{"/127.0.0.1", "GET", 200},
		{"/:ip", "GET", 200},
	}},
}

func myMethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	})
}

func myMethodNotFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	})
}

func myPanicHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "ne ne ne", 500)
	})
}

func TestRouter(t *testing.T) {
	router := New()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)

	router.ServeHTTP(w, req)
	expect(t, w.Code, http.StatusOK)
	expect(t, len(w.HeaderMap), 0)
}

func TestRoutes(t *testing.T) {
	router := New()
	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	for _, v := range routes {
		if len(v.methods) < 1 {
			v.methods = "ALL"
		}
		router.HandleFunc(v.path, func(w http.ResponseWriter, r *http.Request) {}, v.methods)

		var w *httptest.ResponseRecorder

		for _, v := range v.requests {
			w = httptest.NewRecorder()
			req, _ := http.NewRequest(v.method, v.request, nil)
			router.ServeHTTP(w, req)
			expect(t, w.Code, v.expect)
			if w.Code != v.expect {
				log.Fatalf("[%s - %s - %d > %d]", v.request, v.method, v.expect, w.Code)
			}
		}
	}
}

func TestPanic(t *testing.T) {
	router := New()
	router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("si si si")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router.ServeHTTP(w, req)
	expect(t, w.Code, http.StatusInternalServerError)
}

func TestPanicHandler(t *testing.T) {
	router := New()
	router.PanicHandler = myPanicHandler()
	router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("ja ja ja")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router.ServeHTTP(w, req)
	expect(t, w.Code, http.StatusInternalServerError)
	expect(t, w.Body.String(), "ne ne ne\n")
}

func TestHandleFunc(t *testing.T) {
	router := New()
	err := router.HandleFunc("/:none", func(w http.ResponseWriter, r *http.Request) {})
	if err == nil {
		t.Error(err)
	}
	err = router.HandleFunc("/*/test", func(w http.ResponseWriter, r *http.Request) {})
	if err == nil {
		t.Error(err)
	}
	router.HandleFunc("/verbose", func(w http.ResponseWriter, r *http.Request) {})
}

func TestNotAllowedHandler(t *testing.T) {
	router := New()
	router.NotAllowedHandler = myMethodNotAllowed()
	router.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {}, "GET")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 405)
}

func TestNotFoundHandler(t *testing.T) {
	router := New()
	router.NotFoundHandler = myMethodNotFound()
	router.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 404)
}

func TestLogRequests(t *testing.T) {
	router := New()
	router.LogRequests = true
	err := router.HandleFunc("/logrequest", func(w http.ResponseWriter, r *http.Request) {})
	expect(t, err, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/logrequest", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestRequestId(t *testing.T) {
	router := New()
	router.LogRequests = true
	router.RequestID = "Request_log_id"
	err := router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	expect(t, err, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Request_log_id", "5629498000ff0daa102de72aef0001737e7a756e7a756e6369746f2d617069000131000100")
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
	expect(t, w.HeaderMap.Get("Request_log_id"), "5629498000ff0daa102de72aef0001737e7a756e7a756e6369746f2d617069000131000100")
}

func TestRequestIdCreate(t *testing.T) {
	router := New()
	router.LogRequests = true
	router.RequestID = "Request-ID"
	err := router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	expect(t, err, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
	expect(t, len(w.HeaderMap.Get("Request-ID")), 0)
}

func TestHandleFuncMethods(t *testing.T) {
	router := New()

	get_handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I handle GET"))
	}
	post_handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I handle POST"))
	}

	router.HandleFunc("/spine", get_handler, "GET")
	router.HandleFunc("/spine", post_handler, "POST")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/spine", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 405)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/spine", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Body.String(), "I handle GET")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/spine", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Body.String(), "I handle POST")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("HEAD", "/spine", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 405)
}

func TestContextNamedParams(t *testing.T) {
	router := New()

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(ParamsKey).(Params)
		if r.Method == "POST" {
			expect(t, params[":uuid"], "A97F0AF3-043D-4376-82BE-CD6C1A524E0E")
		}
		if r.Method == "GET" {
			expect(t, params["*"], "catch-all-context")
		}
		w.Write([]byte("named params"))
	}

	router.HandleFunc("/test/:uuid", handler, "POST")
	router.HandleFunc("/test/*", handler, "GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test/A97F0AF3-043D-4376-82BE-CD6C1A524E0E", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)

	req, _ = http.NewRequest("GET", "/test/catch-all-context", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestContextMiddleware(t *testing.T) {
	router := New()

	// Test middleware with context
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "m1", "m1")
			ctx = context.WithValue(ctx, "key", 1)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := r.Context().Value(ParamsKey).(Params)
			ctx := context.WithValue(r.Context(), "m2", "m2")
			ctx = context.WithValue(ctx, "uuid val", params[":uuid"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	m3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "m3", "m3")
			ctx = context.WithValue(ctx, "ctx", "string")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(ParamsKey).(Params)
		expect(t, r.Context().Value("m1"), "m1")
		expect(t, r.Context().Value("m2"), "m2")
		expect(t, r.Context().Value("m3"), "m3")
		expect(t, r.Context().Value("uuid val"), "A97F0AF3-043D-4376-82BE-CD6C1A524E0E")
		expect(t, params[":uuid"], "A97F0AF3-043D-4376-82BE-CD6C1A524E0E")
		expect(t, r.Context().Value("ctx"), "string")
		expect(t, r.Context().Value("key"), 1)
		w.Write([]byte("named params"))
	}

	stdChain := middleware.New(m1, m2, m3)
	router.Handle("/foo/:uuid", stdChain.ThenFunc(handler), "PATCH")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/foo/A97F0AF3-043D-4376-82BE-CD6C1A524E0E", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestContextNamedParamsSlice(t *testing.T) {
	router := New()

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(ParamsKey).(Params)
		uuid := params[":uuid"]
		fmt.Printf("uuid = %+v\n", len(uuid.(string)))
		fmt.Printf("uuid = %+v\n", uuid.(string)[0])

		//		expect(t, uuid[0], "A97F0AF3-043D-4376-82BE-CD6C1A524E0E")
		//		expect(t, uuid[1], "12EC2DA8-403D-4C8B-AE39-D011762181A0")
		w.Write([]byte("named params"))
	}

	router.HandleFunc("/test/:uuid/:uuid", handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/A97F0AF3-043D-4376-82BE-CD6C1A524E0E/12EC2DA8-403D-4C8B-AE39-D011762181A0", nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}
