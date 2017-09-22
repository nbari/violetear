// go test -run=BenchmarkRouter -bench=.
// go test -bench=.
// go test -bench=BenchmarkRouter

package test_violetear

import (
	"net/http"
	"testing"

	"github.com/nbari/violetear"
)

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := &violetear.ResponseWriter{}
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

func BenchmarkRouter(b *testing.B) {
	router := violetear.New()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {})
	r, _ := http.NewRequest("GET", "/hello", nil)
	benchRequest(b, router, r)
}
