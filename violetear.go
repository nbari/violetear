/**
 * HTTP Router
 *
 */
package violetear

import (
	"log"
	"net/http"
	//"net/url"
	//"regexp"
	//"strings"
)

type Param struct {
	Key   string
	Value string
}

type Route struct {
	regex   string
	handler http.Handler
	method  [3]string
}

type Router struct {
	NotFoundHandler  http.Handler
	MethodNotAllowed http.HandlerFunc
	routes           []*Route
}

func (r Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := ah.h(ah.appContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
