// go test -run=BenchmarkRouter -bench=.
// go test -bench=.
// go test -bench=BenchmarkRouter

package violetear

import (
	"net/http"
	"testing"
)

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := &ResponseWriter{}
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}

func BenchmarkRouterStatic(b *testing.B) {
	router := New()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {}, "GET,HEAD")
	r, _ := http.NewRequest("GET", "/hello", nil)
	benchRequest(b, router, r)
}

func BenchmarkRouterDynamic(b *testing.B) {
	router := New()
	router.AddRegex(":word", `^\w+$`)
	router.HandleFunc("/test/:word", func(w http.ResponseWriter, r *http.Request) {}, "GET,HEAD")
	r, _ := http.NewRequest("GET", "/test/foo", nil)
	benchRequest(b, router, r)
}
