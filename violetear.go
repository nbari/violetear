// HTTP router
package violetear

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Violetear struct {
	// Routes to be matched
	routes *Trie

	// map of dynamic routes and regular expresions
	dynamicRoutes dynamicSet

	// log requests
	logRequests bool

	// to check
	PanicHandler    func(http.ResponseWriter, *http.Request, interface{})
	NotFoundHandler http.Handler
}

var (
	split_path_rx = regexp.MustCompile(`[^/ ]+`)
)

// New returns a new initialized router.
func New(log ...bool) *Violetear {
	v := &Violetear{
		routes:        NewTrie(),
		dynamicRoutes: make(dynamicSet),
	}

	if len(log) > 0 {
		v.logRequests = true
	}
	return v
}

// Run violetear as an HTTP server.
// The addr string takes the same format as http.ListenAndServe.
func (v *Violetear) Run(addr string) {
	log.Printf("Violetear listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, v))
}

func (v *Violetear) AddPath(path string, handler http.HandlerFunc, http_methods ...string) {
	path_parts := v.splitPath(path)

	// search for dynamic routes
	for _, p := range path_parts {
		if strings.HasPrefix(p, ":") {
			if _, ok := v.dynamicRoutes[p]; !ok {
				log.Fatalf("[%s] not found, need to add it using AddRegex(\"%s\", `your regex`)", p, p)
			}
		}
	}

	// if no methods, accept ALL
	methods := "ALL"
	if len(http_methods) > 0 {
		methods = http_methods[0]
	}

	log.Printf("Adding path: %s, Handler: %s, Methods: %s", path, handler, methods)
	v.routes.Set(path_parts, handler, methods)
}

// AddRegex adds a ":named" regular expression to the dynamicRoutes
func (v *Violetear) AddRegex(name string, regex string) error {
	return v.dynamicRoutes.Set(name, regex)
}

// Match matches registered paths against the request.
func (v *Violetear) Match(req *http.Request) bool {
	return false
}

// func (r *Violetear) Handler(path string, handler http.Handler) {}

// func (r *Violetear) HandlerFunc(path string, handler http.HandlerFunc) {}

// ServerHTTP dispatches the handler registered in the matched path
func (v *Violetear) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if v.logRequests {
		log.Println(req.Method, req.RequestURI)
	}

	split_request := v.splitPath(req.RequestURI)

	r, l := v.routes.Get(split_request)

	var handler http.Handler

	if len(r.Handler) > 0 && l {
		handler = r.Handler["ALL"]
	} else if r.HasRegex {
		for k, _ := range r.Node {
			if strings.HasPrefix(k, ":") {
				handler.ServeHTTP(res, req)
			}
		}
	} else {
		log.Print("Not found")
	}
	handler.ServeHTTP(res, req)

}

// splitPath returns an slice of the path
func (v *Violetear) splitPath(p string) []string {
	path_parts := split_path_rx.FindAllString(p, -1)

	// root (empty slice)
	if len(path_parts) == 0 {
		path_parts = append(path_parts, "/")
	}

	return path_parts
}
