package main

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/nbari/violetear"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	router := violetear.New()
	router.LogRequests = true
	router.HandleFunc("/", hello)
	// Register pprof handlers
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	log.Fatal(http.ListenAndServe(":8080", router))
}
