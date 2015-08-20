/**
 * HTTP Router
 *
 */
package violetear

import (
	"net/http"
	"regexp"
)

type Router struct {
	// API versions
	Versions []string

	// Hosts (host, vroot)
	Hosts map[string]string

	//dynamic (vroot, hosts)
	DynamicHosts map[string]regexp.Regexp
}

func New(file string) Config {
	var config Config
	config = GetConfig(file)
	return config
}

func Add(resource string, handler string, methods []string) Config {
	var config Config
	return config
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	return
}
