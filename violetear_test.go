package violetear

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectDeepEqual(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestRouter(t *testing.T) {
	router := New()
	router.SetHeader("X-app-epazote", "1.1")
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)

	router.ServeHTTP(w, req)
	expect(t, w.Code, http.StatusOK)
	expect(t, len(w.HeaderMap), 2)
	expectDeepEqual(t, w.HeaderMap["X-App-Epazote"], []string{"1.1"})
	fmt.Println(w.Body)
}
