// HTTP router
package violetear

import (
	"log"
	_ "net/http"
	"regexp"
	"strings"
)

type Router struct {
	routes        *Trie
	dynamicRoutes dynamicSet
}

func New() *Router {
	log.Print("Starting violetear ...")

	return &Router{
		routes:        NewTrie(),
		dynamicRoutes: make(dynamicSet),
	}
}

func (r *Router) Add(path string, handler string, http_methods ...string) {

	// create slice from path
	rx := regexp.MustCompile(`[^/ ]+`)
	path_parts := rx.FindAllString(path, -1)
	// root (empty slice)
	if len(path_parts) == 0 {
		path_parts = append(path_parts, "/")
	}

	// search for dynamic routes
	for _, v := range path_parts {
		if strings.HasPrefix(v, ":") {
			if _, ok := r.dynamicRoutes[v]; !ok {
				log.Fatalf("[%s] not found, need to add it using AddRegex(\"%s\", `your regex`)", v, v)
			}
		}
	}

	// if no methods, accept ALL
	methods := "ALL"
	if len(http_methods) > 0 {
		methods = http_methods[0]
	}

	log.Printf("Adding path: %s, Handler: %s, Methods: %s", path, handler, methods)
	r.routes.Set(path_parts, handler, methods)

}

func (r *Router) AddRegex(name string, regex string) error {
	return r.dynamicRoutes.Set(name, regex)
}
