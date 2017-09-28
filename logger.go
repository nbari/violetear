package violetear

import (
	"log"
	"net/http"
)

// logger log values separated by space
func logger(ww *ResponseWriter, r *http.Request) {
	log.Printf("%s [%s] %d %d %s %s",
		r.RemoteAddr,
		r.URL,
		ww.Status(),
		ww.Size(),
		ww.RequestTime(),
		ww.RequestID())
}
