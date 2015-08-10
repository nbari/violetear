/**
 * HTTP Router
 *
 */
package violetear

import (
	"net/http"
	"regexp"
	//"strings"
)

type Router struct {
	// API versions
	Versions []string

	// Hosts (host, vroot)
	Hosts map[string]string

	//dynamic (vroot, hosts)
	DynamicHosts map[string]regexp.Regexp

	StaticRoute map[string]StaticRoute
}

type StaticRoute struct {
	URL, handler string
	Methods      []string
}

//var _ http.Handler = New()

func New(file string) Config {
	var config Config
	config = GetConfig(file)
	return config
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
