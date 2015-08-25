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

	// dynamicRoutes map of dynamic routes and regular expresions
	dynamicRoutes dynamicSet

	// logRequests yes or no
	logRequests bool

	// NotFoundHandler configurable http.Handler which is called when no matching
	// route is found. If it is not set, http.NotFound is used.
	NotFoundHandler http.Handler

	// NotAllowedHandler configurable http.Handler which is called when method not allowed.
	NotAllowedHandler http.Handler

	// Function to handle panics recovered from http handlers.
	PanicHandler func(http.ResponseWriter, *http.Request)
}

var split_path_rx = regexp.MustCompile(`[^/ ]+`)

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

// HandleFunc add a route to the router (path, HandlerFunc, methods)
func (v *Violetear) HandleFunc(path string, handler http.HandlerFunc, http_methods ...string) {
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

// MethodNotAllowed default handler for 405
func (v *Violetear) MethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	})
}

// ServerHTTP dispatches the handler registered in the matched path
func (v *Violetear) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if v.logRequests {
		log.Println(req.Method, req.RequestURI)
	}

	node, path, leaf := v.routes.Get(v.splitPath(req.RequestURI))

	// checkMethod check if method is allowed or not
	checkMethod := func(node *Trie, method string) http.Handler {
		if h, ok := node.Handler[method]; ok {
			return h
		} else if h, ok := node.Handler["ALL"]; ok {
			return h
		}
		if v.NotAllowedHandler != nil {
			return v.NotAllowedHandler
		} else {
			return v.MethodNotAllowed()
		}
	}

	var match func(node *Trie, path []string, leaf bool) http.Handler

	// match find a handler for the request
	match = func(node *Trie, path []string, leaf bool) http.Handler {
		if len(node.Handler) > 0 && leaf {
			return checkMethod(node, req.Method)
		} else if node.HasRegex {
			for k, _ := range node.Node {
				if strings.HasPrefix(k, ":") {
					rx := v.dynamicRoutes[k]
					if rx.MatchString(path[0]) {
						path[0] = k
						if leaf {
							match(node, path, leaf)
						} else {
							return checkMethod(node.Node[k], req.Method)
						}
					}
				}
			}
		}
		if v.NotFoundHandler != nil {
			return v.NotFoundHandler
		}
		return http.NotFoundHandler()
	}

	//var handler http.Handler
	h := match(node, path, leaf)
	log.Printf("%T", h)
	res.Header().Set("X-app-epazote", "1.0")
	h.ServeHTTP(res, req)
	return
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
