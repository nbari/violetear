// HTTP router
//
// Basic example:
//
//  package main
//
//  import (
//     "fmt"
//     "github.com/nbari/violetear"
//     "log"
//     "net/http"
//  )
//
//  func catchAll(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, r.URL.Path[1:])
//  }
//
//  func helloWorld(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, r.URL.Path[1:])
//  }
//
//  func handleUUID(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, r.URL.Path[1:])
//  }
//
//  func main() {
//      router := violetear.New()
//      router.LogRequests = true
//      router.RequestID = "REQUEST_LOG_ID"
//
//      router.AddRegex(":uuid", `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
//
//      router.HandleFunc("*", catchAll)
//      router.HandleFunc("/hello/", helloWorld, "GET,HEAD")
//      router.HandleFunc("/root/:uuid/item", handleUUID, "POST,PUT")
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  }
//
package violetear

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Router struct {
	// Routes to be matched
	routes *Trie

	// dynamicRoutes map of dynamic routes and regular expresions
	dynamicRoutes dynamicSet

	// LogRequests yes or no
	LogRequests bool

	// NotFoundHandler configurable http.Handler which is called when no matching
	// route is found. If it is not set, http.NotFound is used.
	NotFoundHandler http.Handler

	// NotAllowedHandler configurable http.Handler which is called when method not allowed.
	NotAllowedHandler http.Handler

	// PanicHandler function to handle panics.
	PanicHandler http.HandlerFunc

	// RequestID name of the header to use or create.
	RequestID string
}

var split_path_rx = regexp.MustCompile(`[^/ ]+`)

// New returns a new initialized router.
func New() *Router {
	return &Router{
		routes:        NewTrie(),
		dynamicRoutes: make(dynamicSet),
	}
}

// Handle registers the handler for the given pattern (path, http.Handler, methods).
func (v *Router) Handle(path string, handler http.Handler, http_methods ...string) error {
	path_parts := v.splitPath(path)

	// search for dynamic routes
	for _, p := range path_parts {
		if strings.HasPrefix(p, ":") {
			if _, ok := v.dynamicRoutes[p]; !ok {
				return fmt.Errorf("[%s] not found, need to add it using AddRegex(\"%s\", `your regex`)", p, p)
			}
		}
	}

	// if no methods, accept ALL
	methods := "ALL"
	if len(http_methods) > 0 {
		methods = http_methods[0]
	}

	log.Printf("Adding path: %s [%s]", path, methods)

	if err := v.routes.Set(path_parts, handler, methods); err != nil {
		return err
	}
	return nil
}

// HandleFunc add a route to the router (path, http.HandlerFunc, methods)
func (v *Router) HandleFunc(path string, handler http.HandlerFunc, http_methods ...string) error {
	return v.Handle(path, handler, http_methods...)
}

// AddRegex adds a ":named" regular expression to the dynamicRoutes
func (v *Router) AddRegex(name string, regex string) error {
	return v.dynamicRoutes.Set(name, regex)
}

// MethodNotAllowed default handler for 405
func (v *Router) MethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	})
}

// ServerHTTP dispatches the handler registered in the matched path
func (v *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := NewResponseWriter(w)

	// panic handler
	defer func() {
		if err := recover(); err != nil {
			if v.PanicHandler != nil {
				v.PanicHandler(w, r)
			} else {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}
	}()

	lw.Context["sopas"] = "si si si"

	// _ path never empty, defaults to ("/")
	node, path, leaf, _ := v.routes.Get(v.splitPath(r.URL.Path))

	// checkMethod check if method is allowed or not
	checkMethod := func(node *Trie, method string) http.Handler {
		if h, ok := node.Handler[method]; ok {
			return h
		} else if h, ok := node.Handler["ALL"]; ok {
			return h
		}
		if v.NotAllowedHandler != nil {
			return v.NotAllowedHandler
		}
		return v.MethodNotAllowed()
	}

	var match func(node *Trie, path []string, leaf bool) http.Handler

	// match find a handler for the request
	match = func(node *Trie, path []string, leaf bool) http.Handler {
		catchall := false
		if len(node.Handler) > 0 && leaf {
			return checkMethod(node, r.Method)
		} else if node.HasRegex {
			for _, n := range node.Node {
				if strings.HasPrefix(n.path, ":") {
					rx := v.dynamicRoutes[n.path]
					if rx.MatchString(path[0]) {
						path[0] = n.path
						node, path, leaf, _ := node.Get(path)
						return match(node, path, leaf)
					}
				}
			}
			if node.HasCatchall {
				catchall = true
			}
		} else if node.HasCatchall {
			catchall = true
		}
		if catchall {
			for _, n := range node.Node {
				if n.path == "*" {
					return checkMethod(n, r.Method)
				}
			}
		}
		// NotFound
		if v.NotFoundHandler != nil {
			return v.NotFoundHandler
		}
		return http.NotFoundHandler()
	}

	// Request-ID
	if v.RequestID != "" {
		if rid := r.Header.Get(v.RequestID); rid != "" {
			lw.Header().Set(v.RequestID, rid)
		}
	}

	//h http.Handler
	h := match(node, path, leaf)

	// dispatch request
	h.ServeHTTP(lw, r)

	if v.LogRequests {
		log.Printf("%s [%s] %d %d %v %s",
			r.RemoteAddr,
			r.URL,
			lw.Status(),
			lw.Size(),
			time.Since(start),
			lw.Header().Get(v.RequestID))
	}
	return
}

// splitPath returns an slice of the path
func (v *Router) splitPath(p string) []string {
	path_parts := split_path_rx.FindAllString(p, -1)

	// root (empty slice)
	if len(path_parts) == 0 {
		path_parts = append(path_parts, "/")
	}

	return path_parts
}
