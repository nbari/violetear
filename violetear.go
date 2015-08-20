/**
 * HTTP Router
 *
 */
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
	log.Print("Starting...")

	return &Router{
		routes:        NewTrie(),
		dynamicRoutes: NewDynamicSet(),
	}
}

func (r *Router) Add(path string, handler string, http_methods ...string) {

	rx := regexp.MustCompile(`[^/ ]+`)
	path_parts := rx.FindAllString(path, -1)

	methods := "ALL"
	if len(http_methods) > 0 {
		methods = http_methods[0]
	}

	log.Printf("Adding path: %s Handler: %s Methods: %s", path_parts, handler, methods)
	r.routes.Set(path_parts, handler, methods)

}
