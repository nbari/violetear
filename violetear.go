// Package violetear - HTTP router
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
//      router.HandleFunc("/hello", helloWorld, "GET,HEAD")
//      router.HandleFunc("/root/:uuid/item", handleUUID, "POST,PUT")
//
//      srv := &http.Server{
//          Addr:           ":8080",
//          Handler:        router,
//          ReadTimeout:    5 * time.Second,
//          WriteTimeout:   7 * time.Second,
//          MaxHeaderBytes: 1 << 20,
//      }
//      log.Fatal(srv.ListenAndServe())
//  }
//
package violetear

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// ParamsKey used for the context
const (
	ParamsKey     key = 0
	versionHeader     = "application/vnd."
)

// key int is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// Router struct
type Router struct {
	// dynamicRoutes map of dynamic routes and regular expressions
	dynamicRoutes dynamicSet

	// Routes to be matched
	routes *Trie

	// Logger
	Logger func(*ResponseWriter, *http.Request)

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

	// Verbose
	Verbose bool

	// Error resulted from building a route.
	err error
}

// New returns a new initialized router.
func New() *Router {
	return &Router{
		dynamicRoutes: dynamicSet{},
		routes:        &Trie{},
		Logger:        logger,
		Verbose:       true,
	}
}

// Handle registers the handler for the given pattern (path, http.Handler, methods).
func (r *Router) Handle(path string, handler http.Handler, httpMethods ...string) *Trie {
	var version string
	if i := strings.Index(path, "#"); i != -1 {
		version = path[i+1:]
		path = path[:i]
	}
	pathParts := r.splitPath(path)

	// search for dynamic routes
	for _, p := range pathParts {
		if strings.HasPrefix(p, ":") {
			if _, ok := r.dynamicRoutes[p]; !ok {
				r.err = fmt.Errorf("[%s] not found, need to add it using AddRegex(%q, `your regex`", p, p)
				return nil
			}
		}
	}

	// if no methods, accept ALL
	methods := "ALL"
	if len(httpMethods) > 0 && len(strings.TrimSpace(httpMethods[0])) > 0 {
		methods = httpMethods[0]
	}

	if r.Verbose {
		log.Printf("Adding path: %s [%s] %s", path, methods, version)
	}

	trie, err := r.routes.Set(pathParts, handler, methods, version)
	if err != nil {
		r.err = err
		return nil
	}
	return trie
}

// HandleFunc add a route to the router (path, http.HandlerFunc, methods)
func (r *Router) HandleFunc(path string, handler http.HandlerFunc, httpMethods ...string) *Trie {
	return r.Handle(path, handler, httpMethods...)
}

// AddRegex adds a ":named" regular expression to the dynamicRoutes
func (r *Router) AddRegex(name, regex string) error {
	return r.dynamicRoutes.Set(name, regex)
}

// MethodNotAllowed default handler for 405
func (r *Router) MethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	})
}

// checkMethod check if request method is allowed or not
func (r *Router) checkMethod(node *Trie, method string) http.Handler {
	for _, h := range node.Handler {
		if h.Method == "ALL" {
			return h.Handler
		}
		if h.Method == method {
			return h.Handler
		}
	}
	if r.NotAllowedHandler != nil {
		return r.NotAllowedHandler
	}
	return r.MethodNotAllowed()
}

// dispatch request
func (r *Router) dispatch(node *Trie, key, path, method, version string, leaf bool, params Params) (http.Handler, Params) {
	catchall := false
	if node.name != "" {
		if params == nil {
			params = Params{}
		}
		params.Add("rname", node.name)
	}
	if len(node.Handler) > 0 && leaf {
		return r.checkMethod(node, method), params
	} else if node.HasRegex {
		for _, n := range node.Node {
			if strings.HasPrefix(n.path, ":") {
				rx := r.dynamicRoutes[n.path]
				if rx.MatchString(key) {
					// add param to context
					if params == nil {
						params = Params{}
					}
					params.Add(n.path, key)
					node, key, path, leaf := node.Get(n.path+path, version)
					return r.dispatch(node, key, path, method, version, leaf, params)
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
				// add "*" to context
				if params == nil {
					params = Params{}
				}
				params.Add("*", key)
				if n.name != "" {
					params.Add("rname", n.name)
				}
				return r.checkMethod(n, method), params
			}
		}
	}
	// NotFound
	if r.NotFoundHandler != nil {
		return r.NotFoundHandler, params
	}
	return http.NotFoundHandler(), params
}

// ServeHTTP dispatches the handler registered in the matched path
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// panic handler
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %s", err)
			if r.PanicHandler != nil {
				r.PanicHandler(w, req)
			} else {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}
	}()

	// Request-ID
	var rid string
	if r.RequestID != "" {
		if rid = req.Header.Get(r.RequestID); rid != "" {
			w.Header().Set(r.RequestID, rid)
		}
	}

	// wrap ResponseWriter
	var ww *ResponseWriter
	if r.LogRequests {
		ww = NewResponseWriter(w, rid)
	}

	// set version based on the value of "Accept: application/vnd.*"
	version := req.Header.Get("Accept")
	if i := strings.LastIndex(version, versionHeader); i != -1 {
		version = version[len(versionHeader)+i:]
	} else {
		version = ""
	}

	// query the path from left to right
	node, key, path, leaf := r.routes.Get(req.URL.Path, version)

	// dispatch the request
	h, p := r.dispatch(node, key, path, req.Method, version, leaf, nil)

	// dispatch request
	if r.LogRequests {
		if p == nil {
			h.ServeHTTP(ww, req)
		} else {
			h.ServeHTTP(ww, req.WithContext(context.WithValue(req.Context(), ParamsKey, p)))
		}
		r.Logger(ww, req)
	} else {
		if p == nil {
			h.ServeHTTP(w, req)
		} else {
			h.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), ParamsKey, p)))
		}
	}
}

// splitPath returns an slice of the path
func (r *Router) splitPath(p string) []string {
	pathParts := strings.FieldsFunc(p, func(c rune) bool {
		return c == '/'
	})
	// root (empty slice)
	if len(pathParts) == 0 {
		return []string{"/"}
	}
	return pathParts
}

// GetError returns an error resulted from building a route, if any.
func (r *Router) GetError() error {
	return r.err
}
