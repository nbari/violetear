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

func (v *Violetear) abc(route *Trie, path []string, leaf bool) http.Handler {
	log.Print(route, path, leaf)
	if len(route.Handler) > 0 && leaf {
		return route.Handler["ALL"]
	} else if route.HasRegex {
		for k, _ := range route.Node {
			if strings.HasPrefix(k, ":") {
				rx := v.dynamicRoutes[k]
				if rx.MatchString(path[0]) {
					log.Print(path, "-------<<<<")
					path[0] = k
					if leaf {
						v.abc(route, path, leaf)
					} else {
						return route.Node[k].Handler["ALL"]
					}
				}
			}
			return route.Handler["ALL"]
		}
	} else {
		return nil
	}
	return route.Handler["ALL"]

}

// ServerHTTP dispatches the handler registered in the matched path
func (v *Violetear) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if v.logRequests {
		log.Println(req.Method, req.RequestURI, req.URL, req.URL.Path)
	}

	split_request := v.splitPath(req.RequestURI)

	route, path, leaf := v.routes.Get(split_request)

	var handler http.Handler

	handler = v.abc(route, path, leaf)
	handler.ServeHTTP(res, req)

	/*
		if len(route.Handler) > 0 && leaf {
			handler = route.Handler["ALL"]
		} else if route.HasRegex {
			for k, _ := range route.Node {
				if strings.HasPrefix(k, ":") {
					rx := v.dynamicRoutes[k]
					if rx.MatchString(path[0]) {
						log.Print(path, "-------<<<<-")
					}
					handler.ServeHTTP(res, req)
				}
			}
			handler = nil
		} else {
			log.Print("Not found")
			handler = nil
		}
	*/

	//handler.ServeHTTP(res, req)

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
