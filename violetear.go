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

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when method not allowed
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	PanicHandler func(http.ResponseWriter, *http.Request)
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

	log.Printf("Adding path: %s, Handler: %T, Methods: %s", path, handler, methods)
	v.routes.Set(path_parts, handler, methods)
}

// AddRegex adds a ":named" regular expression to the dynamicRoutes
func (v *Violetear) AddRegex(name string, regex string) error {
	return v.dynamicRoutes.Set(name, regex)
}

// func (r *Violetear) Handler(path string, handler http.Handler) {}

// func (r *Violetear) HandlerFunc(path string, handler http.HandlerFunc) {}

// Match matches registered paths against the request.
func (v *Violetear) Match(node *Trie, path []string, leaf bool) map[string]http.Handler {

	log.Print(node, path, leaf)

	if len(node.Handler) > 0 && leaf {
		log.Print("Matched -------- primer round")
		return node.Handler
	} else if node.HasRegex {
		for k, _ := range node.Node {
			if strings.HasPrefix(k, ":") {
				rx := v.dynamicRoutes[k]
				log.Print(path, "trying to find a match -------<<<<")
				if rx.MatchString(path[0]) {
					log.Print(path, "matched -------<<<<")
					path[0] = k
					if leaf {
						v.Match(node, path, leaf)
					} else {
						log.Print("matched regex no leaf, returnint")
						return node.Node[k].Handler
					}
				}
			}
		}
		log.Print("Not found ---------------------")
		return nil
	} else {
		log.Print("Not match ---------------------")
		return nil
	}
}

// ServerHTTP dispatches the handler registered in the matched path
func (v *Violetear) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if v.logRequests {
		log.Println(req.Method, req.RequestURI)
	}

	split_request := v.splitPath(req.RequestURI)

	node, path, leaf := v.routes.Get(split_request)

	var handler http.Handler

	handlers := v.Match(node, path, leaf)
	handler = handlers["ALL"]
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
